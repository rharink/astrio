package server

//Room maintains the set of active connections and broadcasts messages to the
//connections.
type Room struct {
	//identifier
	id string
	//registered connections.
	connections map[*Connection]bool
	//inbound messages from the connections.
	broadcast chan []byte
	//register requests from the connections.
	register chan *Connection
	//unregister requests from connections.
	unregister chan *Connection
}

//create and run a new Room
func NewRoom(id string) (r *Room) {
	r = &Room{
		id:          id,
		broadcast:   make(chan []byte),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
	}

	go r.listen()

	return r
}

//listen for incomming messages
func (r *Room) listen() {
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

//broadcast a message to all users in the Room
func (r *Room) Broadcast(message []byte) {
	r.broadcast <- message
}
