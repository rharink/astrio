package server

//hub maintains the set of active connections and broadcasts messages to the
//connections.
type room struct {
	//identifier
	id string

	//registered connections.
	connections map[*connection]bool

	//inbound messages from the connections.
	broadcast chan []byte

	//register requests from the connections.
	register chan *connection

	//unregister requests from connections.
	unregister chan *connection
}

//create and run a new hub
func newRoom(id string) (r *room) {
	r = &room{
		id:          id,
		broadcast:   make(chan []byte),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
	}

	go r.listen()

	return r
}

func (r *room) listen() {
	for {
		select {
		case c := <-r.register:
			r.connections[c] = true
		case c := <-r.unregister:
			if _, ok := r.connections[c]; ok {
				delete(r.connections, c)
				close(c.send)
			}
		case m := <-r.broadcast:
			for c := range r.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(r.connections, c)
				}
			}
		}
	}
}

func (r *room) Send(message []byte) {
	r.broadcast <- message
}
