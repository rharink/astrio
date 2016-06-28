package game

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/rauwekost/astrio/configuration"
	"github.com/rauwekost/astrio/game/packet"
)

//Connection ...
type Player struct {
	ws     *websocket.Conn
	sendch chan []byte
	game   *Game
}

//NewPlayer returns a new player
func NewPlayer(ws *websocket.Conn, game *Game) *Player {
	return &Player{
		ws:     ws,
		sendch: make(chan []byte, 256),
		game:   game,
	}
}

//Reader pumps messages from the websocket connection to the hub.
func (p *Player) ReadPump() {
	defer func() {
		p.Close()
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
			p.Close()
		}

		if packet.OPCode(message[0]) == packet.OPCodeServerMouseMove {
			pack := packet.MouseMove{}
			packet.Decode(message, &pack)
			fmt.Printf("%+v", pack)
		}
		//send messages to the packet handler
		//ph.OnMessage(message)
	}
}

//Writer pumps messages from the hub to the websocket connection.
func (p *Player) WritePump() {
	pong := time.NewTicker(configuration.Server.PingPeriod)
	defer func() {
		pong.Stop()
		p.Close()
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

//Close closes connections and unregisters the player
func (p *Player) Close() {
	p.game.unregister <- p
	p.ws.Close()
}

//Join player joins the given game
func (p *Player) Join() {
	fmt.Println("player joined")
	p.game.register <- p
	go p.ReadPump()
	go p.WritePump()
}

//wite writes a message with the given message type and payload.
func (p *Player) write(mt int, payload []byte) error {
	p.ws.SetWriteDeadline(time.Now().Add(configuration.Server.WriteWait))

	return p.ws.WriteMessage(mt, payload)
}
