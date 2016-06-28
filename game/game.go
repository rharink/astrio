package game

import "time"

//Game ...
type Game struct {
	//unique identifier for the game
	id string
	//if the game is runninga
	running bool
	//ticker every 40ms (25fps)
	mainLoop *time.Ticker
	//active players in the game
	players map[*Player]bool
	//inbound messages from the connections.
	broadcast chan []byte
	//register requests from the connections.
	register chan *Player
	//unregister requests from connections.
	unregister chan *Player
}

//New return  a new game
func New(id string) *Game {
	return &Game{
		id:         id,
		running:    false,
		mainLoop:   time.NewTicker(40 * time.Millisecond),
		players:    make(map[*Player]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Player),
		unregister: make(chan *Player),
	}
}

//Run run the game
func (g *Game) Run() {
	if g.running == true {
		return
	}

	go func() {
		for range g.mainLoop.C {
		}
	}()

	go g.listen()

	g.running = true
}

//Stop stops the game
func (g *Game) Stop() {
	if g.running == false {
		return
	}

	g.running = false
}

//listen for incomming messages from players
func (g *Game) listen() {
	for {
		select {
		case p := <-g.register:
			g.players[p] = true
		case p := <-g.unregister:
			if _, ok := g.players[p]; ok {
				delete(g.players, p)
				close(p.sendch)
			}
		case m := <-g.broadcast:
			for p := range g.players {
				select {
				case p.sendch <- m:
				default:
					close(p.sendch)
					delete(g.players, p)
				}
			}
		}
	}
}

//broadcast a message to all players in the game
func (g *Game) Broadcast(message []byte) {
	g.broadcast <- message
}
