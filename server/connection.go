package server

import (
	"time"

	"github.com/gorilla/websocket"
)

type connection struct {
	//websocket connection
	ws *websocket.Conn

	//send buffered channel for outbound messages
	send chan []byte

	//pointer to the hub this connection belongs to
	room *room

	//tick time
	pingPeriod time.Duration

	//write wait duration
	writeWait time.Duration
}

func NewConnection(ws *websocket.Conn, room *room, tick time.Duration, writeWait time.Duration) *connection {
	conn := connection{
		ws:         ws,
		send:       make(chan []byte, 256),
		room:       room,
		pingPeriod: tick,
		writeWait:  writeWait,
	}

	go conn.readPump()
	go conn.writePump()

	return &conn
}

//readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		c.room.unregister <- c
		c.ws.Close()
	}()
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.room.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(c.pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

//wite writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(c.writeWait))
	return c.ws.WriteMessage(mt, payload)
}
