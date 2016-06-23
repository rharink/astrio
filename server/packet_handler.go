package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type (
	PacketHandler struct {
		Player          *Player
		Protocol        int
		HandshakePassed bool
		PressQ          bool
		PressW          bool
		PressSpace      bool
	}
)

const (
	UPDATE_MOVEMENT uint8 = 16
)

func NewPacketHandler(p *Player) *PacketHandler {
	return &PacketHandler{
		Player:     p,
		Protocol:   0,
		PressQ:     false,
		PressW:     false,
		PressSpace: false,
	}
}

func (h *PacketHandler) OnMessage(message []byte) {
	buf := bytes.NewReader(message)
	var opcode uint8
	var protocol uint32
	binary.Read(buf, binary.LittleEndian, &opcode)

	if !h.HandshakePassed {
		if opcode != 254 {
			return //wait for handshake
		}
		binary.Read(buf, binary.LittleEndian, &protocol)
		if protocol < 1 || protocol > 8 {
			h.Player.Unregister() //unsupported protocol
		}
	}

	switch opcode {
	case UPDATE_MOVEMENT:
		h.updateMovement(buf)
	default:
		fmt.Println(message)
	}
}

func (h *PacketHandler) updateMovement(buf io.ReadSeeker) {
	var y, x, id float32
	binary.Read(buf, binary.LittleEndian, &x)
	binary.Read(buf, binary.LittleEndian, &y)
	binary.Read(buf, binary.LittleEndian, &id)

	h.Player.Tracker.Mouse.X = x
	h.Player.Tracker.Mouse.Y = y
}
