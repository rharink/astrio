package server

import (
	"fmt"
	"time"

	"github.com/rauwekost/astrio/server/packet"
)

type (
	Game struct {
		//unique identifier for the game
		id string

		//if the game is running
		running bool

		//map of active players holding connections
		players map[*Player]bool

		//inbound messages from the players.
		broadcast chan []byte

		//register player
		register chan *Player

		//unregister Player
		unregister chan *Player

		ticker *time.Ticker
	}
)

//NewGame returns a new Game instance
func NewGame(id string) *Game {
	g := &Game{
		id:         id,
		running:    false,
		players:    make(map[*Player]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		ticker:     time.NewTicker(40 * time.Millisecond),
	}
	go func() {
		for range g.ticker.C {
			u := &packet.UpdateNodes{}
			fmt.Println(u.Bytes())
			g.Send(u.Bytes())
		}
	}()
	go g.listen()

	return g
}

//listen for messages
func (g *Game) listen() {
	for {
		select {
		case p := <-g.register:
			g.players[p] = true
		case p := <-g.unregister:
			g.RemovePlayer(p)
		case m := <-g.broadcast:
			for p := range g.players {
				select {
				case p.sendch <- m:
				default:
					g.RemovePlayer(p)
				}
			}
		}
	}
}

//RemovePlayer removes a player from the game
func (g *Game) RemovePlayer(p *Player) {
	if _, ok := g.players[p]; ok {
		delete(g.players, p)
		close(p.sendch)
	}
}

//send a message
func (g *Game) Send(message []byte) {
	g.broadcast <- message
}
