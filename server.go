package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bmizerany/pat"
	_ "github.com/coreos/dex/pkg/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/websocket"
	"github.com/justinas/alice"
	"github.com/rauwekost/jwt-middleware"
)

type Server struct {
	//address to bind on
	addr string

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

	//size of the upgrader read buffer
	readBufferSize int

	//size of the upgrader write buffer
	writeBufferSize int

	//algorithm that jwt uses
	jwtAlgorithm string

	//path to the jwt key
	jwtKey string

	//httpServer for handling socket transport
	httpServer *http.Server

	//middleware for each request
	middleware alice.Chain
}

//NewServer returns a new server instance based on configuration
func NewServer(cfg *Config) (*Server, error) {
	s := Server{
		addr:            cfg.ServerAddress,
		writeWait:       cfg.ServerWriteWait,
		pongWait:        cfg.ServerPongWait,
		pingPeriod:      cfg.ServerPingPeriod,
		maxMessageSize:  cfg.ServerMaxMessageSize,
		readBufferSize:  cfg.ServerReadBufferSize,
		writeBufferSize: cfg.ServerWriteBufferSize,
		jwtAlgorithm:    cfg.ServerJWTAlgorithm,
		jwtKey:          cfg.ServerJWTKey,
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
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(ws)
}

//init initializes the server befor running
func (s *Server) init() error {
	log.Infof("initializing...")
	log.Infof("version [%s] build[%s] buildDate[%s]", version, build, date)

	//middleware
	log.Info("creating middleware...")
	m := jwtmiddleware.Middleware{
		ParameterName: "token",
		Keyfunc: func(t *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		Successfunc: func(r *http.Request, t *jwt.Token) {
			fmt.Println("Hoera jwt is oke!")
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
			return true
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

	return nil
}
