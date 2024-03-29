package game

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
		// Create a new game
		game = New(key)
		h.games[key] = game
	}

	return game
}
