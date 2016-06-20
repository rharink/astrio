package main

import "sync"

//Hub holds multiple rooms
type Hub struct {
	sync.RWMutex
	rooms map[string]*room
}

//NewHub returns a new Hub
func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*room, 0),
	}
}

//Get a room by key
func (h *Hub) Get(key string) *room {
	h.RLock()
	defer h.RUnlock()

	room, ok := h.rooms[key]

	if !ok {
		// Create a new room
		room = newRoom(key)
		h.rooms[key] = room
	}

	return room
}
