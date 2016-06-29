package game

import (
	"fmt"
	"time"

	cfg "github.com/rauwekost/astrio/configuration"
	"github.com/rauwekost/astrio/game/packet"
)

//Game ...
type Game struct {
	//unique identifier for the game
	id string
	//if the game is runninga
	running bool
	//ticker every 40ms (25fps)
	ticker *time.Ticker
	//the last tick
	lastTick time.Time
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

	//listen to incomming registers, unregisters and broadcasts.
	go g.listen()

	//mainloop ticker trigger
	go func() {
		for t := range g.ticker.C {
			g.lastTick = t
			g.mainLoop()
		}
	}()

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
	fmt.Println("tick", g.lastTick)
	update := packet.UpdateNodes{}
	for p, _ := range g.players {
		update.Nodes = append(update.Nodes, packet.Node{
			ID:    1,
			X:     p.Tracker.Mouse.X,
			Y:     p.Tracker.Mouse.Y,
			Size:  40,
			R:     255,
			G:     255,
			B:     255,
			Flags: 1,
		})
	}
	b := update.Bytes()
	fmt.Println(b)
	g.Broadcast(b)
}

//broadcast a message to all players in the game
func (g *Game) Broadcast(message []byte) {
	g.broadcast <- message
}
