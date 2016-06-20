package server

type (
	Game struct {
		//unique identifier for the game
		id string

		//if the game is running
		running bool

		//map of active connections
		connections map[*connection]bool

		//inbound messages from the connections.
		broadcast chan []byte

		//register requests from the connections.
		register chan *connection

		//unregister requests from connections.
		unregister chan *connection
	}
)

func NewGame(id string) *Game {
	g := &Game{
		id:          id,
		running:     false,
		connections: make(map[*connection]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}

	go g.listen()

	return g
}

//listen for messages
func (g *Game) listen() {
	for {
		select {
		case c := <-g.register:
			g.connections[c] = true
		case c := <-g.unregister:
			if _, ok := g.connections[c]; ok {
				delete(g.connections, c)
				close(c.send)
			}
		case m := <-g.broadcast:
			for c := range g.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(g.connections, c)
				}
			}
		}
	}
}

//send a message
func (g *Game) Send(message []byte) {
	g.broadcast <- message
}
