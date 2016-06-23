package server

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/rauwekost/astrio/configuration"
)

//Player ...
type Player struct {
	ID            string
	Name          string
	ws            *websocket.Conn
	Game          *Game
	sendch        chan []byte
	Tracker       *PlayerTracker
	packetHandler *PacketHandler
}

//NewPlayer returns a new player
func NewPlayer(ws *websocket.Conn, game *Game, id, name string) *Player {
	p := Player{
		ID:      id,
		Name:    name,
		sendch:  make(chan []byte, 256),
		ws:      ws,
		Game:    game,
		Tracker: &PlayerTracker{},
	}

	p.packetHandler = NewPacketHandler(&p)
	return &p
}

//Reader pumps messages from the websocket connection to the hub.
func (p *Player) Reader() {
	defer func() {
		p.Unregister()
	}()
	for {
		_, message, err := p.ws.ReadMessage()
		if err != nil {
			break
		}

		//ignore empty messages
		if len(message) == 0 {
			logrus.Warn("empty message skipped")
			continue
		}

		//spam
		if len(message) > 2048 {
			logrus.Warn("player unregistered because of spam")
			p.Unregister()
		}

		p.packetHandler.OnMessage(message)
		p.Game.broadcast <- message
	}
}

//Writer pumps messages from the hub to the websocket connection.
func (p *Player) Writer() {
	pong := time.NewTicker(configuration.Server.PingPeriod)
	defer func() {
		pong.Stop()
		p.ws.Close()
	}()
	for {
		select {
		case message, ok := <-p.sendch:
			if !ok {
				p.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := p.write(websocket.BinaryMessage, message); err != nil {
				return
			}
		case <-pong.C:
			if err := p.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

//Unregister unregister a player from a game and close the socket
func (p *Player) Unregister() {
	p.Game.unregister <- p
	p.ws.Close()
}

//wite writes a message with the given message type and payload.
func (p *Player) write(mt int, payload []byte) error {
	p.ws.SetWriteDeadline(time.Now().Add(configuration.Server.WriteWait))

	return p.ws.WriteMessage(mt, payload)
}
