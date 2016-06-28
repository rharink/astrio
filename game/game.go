package game

import (
	"fmt"
	"time"

	cfg "github.com/rauwekost/astrio/configuration"
)

//Game ...
type Game struct {
	//unique identifier for the game
	id string
	//if the game is runninga
	running bool
	//ticker every 40ms (25fps)
	ticker *time.Ticker
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
		ticker:     time.NewTicker(time.Duration(cfg.Game.Tick) * time.Millisecond),
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

	//listen to incomming registers, unregisters, ticks etc.
	go g.listen()

	g.running = true
}

//Stop stops the game
func (g *Game) Stop() {
	if g.running == false {
		return
	}

	g.ticker.Stop()
	g.running = false
}

//listen for incomming messages from players
func (g *Game) listen() {
	for {
		select {
		case <-g.ticker.C:
			g.mainLoop()
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

//mainLoop is the games main loop. this gets called n times a second sending
//updates to all players in the game
func (g *Game) mainLoop() {
	for p, _ := range g.players {
		fmt.Println(p.Tracker)
	}
}

//broadcast a message to all players in the game
func (g *Game) Broadcast(message []byte) {
	g.broadcast <- message
}
