package server

import "sync"

//Hub holds multiple rooms
type Hub struct {
	sync.RWMutex
	rooms map[string]*Room
}

//NewHub returns a new Hub
func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room, 0),
	}
}

//Get a room by key
func (h *Hub) Get(key string) *Room {
	h.RLock()
	defer h.RUnlock()

	room, ok := h.rooms[key]

	if !ok {
		// Create a new room
		room = NewRoom(key)
		h.rooms[key] = room
	}

	return room
}

//Length retuns the number of games
func (h *Hub) Length() int {
	return len(h.rooms)
}

//Stats returns a stats map of games and player amount
func (h *Hub) Stats() map[string]int {
	s := map[string]int{}
	for _, r := range h.rooms {
		s[r.id] = len(r.connections)
	}

	return s
}
