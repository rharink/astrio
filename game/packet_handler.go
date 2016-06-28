package game

import (
	"fmt"

	"github.com/rauwekost/astrio/game/packet"
)

type (
	PacketHandler struct {
		Protocol        int
		HandshakePassed bool
		PressQ          bool
		PressW          bool
		PressSpace      bool
	}
)

//NewPacketHandler returns a new packet handler
func NewPacketHandler() *PacketHandler {
	return &PacketHandler{
		Protocol:        0,
		HandshakePassed: false,
		PressQ:          false,
		PressW:          false,
		PressSpace:      false,
	}
}

//OnMessage when a message comes in
func (h *PacketHandler) OnMessage(m []byte) error {
	switch packet.OPCode(m[0]) {
	case packet.OPCodeServerMouseMove:
		p := packet.MouseMove{}
		err := packet.Decode(m, &p)
		fmt.Println(p)
		return err
	default:
		return nil
	}
}
