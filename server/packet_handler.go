package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type (
	PacketHandler struct {
		Player     *Player
		Protocol   int
		PressQ     bool
		PressW     bool
		PressSpace bool
	}
)

const (
	UPDATE_MOVEMENT uint8 = 16
)

func NewPacketHandler(p *Player) *PacketHandler {
	return &PacketHandler{
		Player:     p,
		Protocol:   1,
		PressQ:     false,
		PressW:     false,
		PressSpace: false,
	}
}

func (h *PacketHandler) OnMessage(message []byte) {
	buf := bytes.NewReader(message)
	var opcode uint8
	binary.Read(buf, binary.LittleEndian, &opcode)

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
