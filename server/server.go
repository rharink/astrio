package server

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/bmizerany/pat"
	_ "github.com/coreos/dex/pkg/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/websocket"
	"github.com/justinas/alice"
	"github.com/rauwekost/astrio/configuration"
	"github.com/rauwekost/jwt-middleware"
)

type Server struct {
	//address to bind on
	addr string

	//frames per second
	fps int

	//time allowed to write a message to the peer
	writeWait time.Duration

	//time allowed to read the nex pong message from the peer
	pongWait time.Duration

	//send pings to peer with this period. Must be less than pongWait
	pingPeriod time.Duration

	//maximum message size allowed from peer
	maxMessageSize int64

	//websocket upgrader
	upgrader websocket.Upgrader

	//hub handles multiple connections comming and going
	hub *Hub

	//size of the upgrader read buffer
	readBufferSize int

	//size of the upgrader write buffer
	writeBufferSize int

	//algorithm that jwt uses
	jwtAlgorithm string

	//path to the jwt key
	jwtSecret string

	//private jwt key
	jwtPrivate string

	//public jwt key
	jwtPublic string

	//allowed orgings
	allowedOrigins []string

	//httpServer for handling socket transport
	httpServer *http.Server

	//middleware for each request
	middleware alice.Chain
}

var (
	Version string
	Build   string
	Date    string
	log     = logrus.WithFields(logrus.Fields{"package": "server"})
)

//NewServer returns a new server instance based on configuration
func New(cfg *configuration.Config) (*Server, error) {
	s := Server{
		addr:            cfg.ServerAddress,
		writeWait:       cfg.ServerWriteWait,
		pongWait:        cfg.ServerPongWait,
		pingPeriod:      cfg.ServerPingPeriod,
		maxMessageSize:  cfg.ServerMaxMessageSize,
		readBufferSize:  cfg.ServerReadBufferSize,
		writeBufferSize: cfg.ServerWriteBufferSize,
		jwtAlgorithm:    cfg.ServerJWTAlgorithm,
		jwtSecret:       cfg.ServerJWTSecret,
		jwtPrivate:      cfg.ServerJWTPrivate,
		jwtPublic:       cfg.ServerJWTPublic,
		allowedOrigins:  cfg.ServerAllowedOrigins,
	}

	return &s, nil
}

//Run runs the server
func (s *Server) Run() error {
	if err := s.init(); err != nil {
		return err
	}

	return s.httpServer.ListenAndServe()
}

//HttpHandler handles http traffic
func (s *Server) httpHandler() http.Handler {
	mux := pat.New()
	mux.Add("GET", "/", s.middleware.ThenFunc(s.handleWebsocket))

	return mux
}

//handleWebsocket handles incomming websocket connections
func (s *Server) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	token := context.Get(r, "token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("forbidden"))
		log.Errorf("upgrader: %s", err)
		return
	}

	ws.SetReadLimit(s.maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(s.pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(s.pongWait)); return nil })
	room := s.hub.Get(claims["room"].(string))
	c := &connection{
		send:       make(chan []byte, 256),
		ws:         ws,
		room:       room,
		pingPeriod: s.pingPeriod,
		writeWait:  s.writeWait,
	}

	room.register <- c
	go c.writePump()
	go c.readPump()
}

func (s *Server) createTempJWT() (string, error) {
	t := jwt.NewWithClaims(signingMethodFromString(s.jwtAlgorithm), jwt.MapClaims{
		"user_id":  1,
		"room":     "astio",
		"team":     "astrio",
		"skin_url": "",
		"exp":      time.Now().Add(300 * time.Second),
	})

	switch signingMethodFromString(s.jwtAlgorithm) {
	case jwt.SigningMethodRS256:
		b, err := ioutil.ReadFile(s.jwtPrivate)
		if err != nil {
			return "", err
		}
		signKey, err := jwt.ParseRSAPrivateKeyFromPEM(b)
		if err != nil {
			return "", err
		}
		return t.SignedString(signKey)
	default:
		return t.SignedString([]byte(s.jwtSecret))
	}
}

//init initializes the server befor running
func (s *Server) init() error {
	log.Infof("initializing...")
	log.Infof("version [%s] build[%s] buildDate[%s]", Version, Build, Date)

	//middleware
	log.Info("creating middleware...")
	m := jwtmiddleware.Middleware{
		ParameterName: "token",
		Keyfunc:       s.getKeyFunc,
		Successfunc: func(r *http.Request, t *jwt.Token) {
			context.Set(r, "token", t)
		},
		Errorfunc: func(err error) {
			log.Error(err)
		},
	}
	s.middleware = alice.New(m.Handler, context.ClearHandler)
	log.Info("middleware created.")

	//upgrader @TODO: make CheckOrigin dynamic through configuration
	log.Info("creating http upgrader...")
	s.upgrader = websocket.Upgrader{
		ReadBufferSize:  s.readBufferSize,
		WriteBufferSize: s.writeBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			for _, o := range s.allowedOrigins {
				if o == "*" || o == r.Header.Get("Origin") {
					return true
				}
			}
			return false
		},
	}
	log.Info("upgrader created.")

	//http server
	log.Info("creating http-server...")
	s.httpServer = &http.Server{
		Addr:    s.addr,
		Handler: s.httpHandler(),
	}
	log.Info("http-server created.")

	//hub
	log.Info("creating hub...")
	s.hub = NewHub()
	log.Info("hub created.")

	//temporary jwt token
	log.Info("creating temporary jwt token")
	t, err := s.createTempJWT()
	if err != nil {
		log.Errorf("error while creating temp jwt-token: %s", err)
	}
	log.Infof("temporary token: %s", t)

	return nil
}

func (s *Server) getKeyFunc(t *jwt.Token) (interface{}, error) {
	switch t.Method {
	case jwt.SigningMethodRS256:
		b, err := ioutil.ReadFile(s.jwtPublic)
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM(b)
		if err != nil {
			return nil, err
		}

		return key, nil
	case jwt.SigningMethodHS256:
		fallthrough
	default:
		return []byte(s.jwtSecret), nil
	}
}

func signingMethodFromString(str string) jwt.SigningMethod {
	switch str {
	case "HS256":
		return jwt.SigningMethodHS256
	case "RS256":
		return jwt.SigningMethodRS256
	default:
		log.Fatalf("unsupported signing-method: %s", str)
		return jwt.SigningMethodHS256
	}
}
