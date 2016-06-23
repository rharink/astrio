package server

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/rauwekost/astrio/configuration"
)

//Connection ...
type Connection struct {
	ws   *websocket.Conn
	send chan []byte
	room *Room
}

//NewConnection returns a new connection
func NewConnection(ws *websocket.Conn, room *Room) *Connection {
	return &Connection{
		ws:   ws,
		send: make(chan []byte, 256),
		room: room,
	}
}

//Reader pumps messages from the websocket connection to the hub.
func (c *Connection) ReadPump() {
	defer func() {
		c.Close()
	}()
	for {
		_, message, err := c.ws.ReadMessage()
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
			logrus.Warn("connection unregistered because of spam")
			c.Close()
		}

		//send messages to the packet handler
		//ph.OnMessage(message)
	}
}

//Writer pumps messages from the hub to the websocket connection.
func (c *Connection) WritePump() {
	pong := time.NewTicker(configuration.Server.PingPeriod)
	defer func() {
		pong.Stop()
		c.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.BinaryMessage, message); err != nil {
				return
			}
		case <-pong.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

//Close closes connection and calls close hook
func (c *Connection) Close() {
	c.room.unregister <- c
	c.ws.Close()
}

//wite writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(configuration.Server.WriteWait))

	return c.ws.WriteMessage(mt, payload)
}
