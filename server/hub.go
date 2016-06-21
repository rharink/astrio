package server

import "sync"

//Hub holds multiple rooms
type Hub struct {
	sync.RWMutex
	games map[string]*Game
}

//NewHub returns a new Hub
func NewHub() *Hub {
	return &Hub{
		games: make(map[string]*Game, 0),
	}
}

//Get a room by key
func (h *Hub) Get(key string) *Game {
	h.RLock()
	defer h.RUnlock()

	game, ok := h.games[key]

	if !ok {
		// Create a new room
		game = NewGame(key)
		h.games[key] = game
	}

	return game
}

//Length retuns the number of games
func (h *Hub) Length() int {
	return len(h.games)
}

//Stats returns a stats map of games and player amount
func (h *Hub) Stats() map[string]int {
	s := map[string]int{}
	for _, g := range h.games {
		s[g.id] = len(g.players)
	}

	return s
}
