package server

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/bmizerany/pat"
	_ "github.com/coreos/dex/pkg/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/websocket"
	"github.com/justinas/alice"
	cfg "github.com/rauwekost/astrio/configuration"
	"github.com/rauwekost/jwt-middleware"
)

var (
	Version string
	Build   string
	Date    string
	log     = logrus.WithFields(logrus.Fields{"package": "server"})
)

//Server a server
type Server struct {
	//websocket upgrader
	upgrader websocket.Upgrader

	//hub handles multiple rooms
	hub *Hub

	//httpServer for handling socket transport
	httpServer *http.Server

	//middleware for each request
	middleware alice.Chain
}

//NewServer returns a new server instance based on cfg
func New() *Server {
	return &Server{
		hub: NewHub(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  cfg.Server.ReadBufferSize,
			WriteBufferSize: cfg.Server.WriteBufferSize,
			CheckOrigin:     checkOrigins,
		},
	}
}

//Run runs the server
func (s *Server) Run() error {
	if err := s.init(); err != nil {
		return err
	}

	log.Infof("server listening on: %s", cfg.Server.Address)
	return s.httpServer.ListenAndServe()
}

//HttpHandler handles http traffic
func (s *Server) httpHandler() http.Handler {
	mux := pat.New()
	mux.Add("GET", "/", s.middleware.ThenFunc(s.handleNewConnection))

	return mux
}

//handleWebsocket handles incomming websocket connections
func (s *Server) handleNewConnection(w http.ResponseWriter, r *http.Request) {
	token := context.Get(r, "token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	ws, err := s.upgradeWebsocket(w, r)
	if err != nil {
		NewForbiddenError(err).Render(w)
		log.Errorf("upgrader: %s", err)
		return
	}

	room := s.hub.Get(claims["game"].(string))
	conn := NewConnection(ws, room)

	room.register <- conn

	go conn.ReadPump()
	go conn.WritePump()
}

//upgrade websocket
func (s *Server) upgradeWebsocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	ws.SetReadLimit(cfg.Server.MaxMessageSize)
	ws.SetReadDeadline(time.Now().Add(cfg.Server.PongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(cfg.Server.PongWait)); return nil })

	return ws, nil
}

//init initializes the server befor running
func (s *Server) init() error {
	log.Infof("initializing...")
	log.Infof("version [%s] build[%s] buildDate[%s]", Version, Build, Date)

	//middleware
	log.Info("creating middleware...")
	m := jwtmiddleware.Middleware{
		ParameterName: "token",
		Keyfunc:       s.getJWTKey,
		Successfunc: func(r *http.Request, t *jwt.Token) {
			context.Set(r, "token", t)
		},
		Errorfunc: func(err error) {
			log.Error(err)
		},
	}
	s.middleware = alice.New(m.Handler, context.ClearHandler)
	log.Info("middleware created.")

	//http transport
	log.Info("creating http server...")
	s.httpServer = &http.Server{
		Addr:    cfg.Server.Address,
		Handler: s.httpHandler(),
	}
	log.Info("http server created.")

	//temporary jwt token
	log.Info("creating temporary jwt token")
	t, err := s.createJWT(&jwt.MapClaims{
		"user": 1,
		"game": "astrio",
		"team": "astrio",
		"exp":  time.Now().Add(300 * time.Second),
	})
	if err != nil {
		log.Errorf("error while creating temp jwt-token: %s", err)
	}
	log.Infof("temporary token: %s", t)

	return nil
}
