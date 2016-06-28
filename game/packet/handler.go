package packet

import (
	"fmt"
)

type Handler struct {
	Protocol        int
	HandshakePassed bool
	PressQ          bool
	PressW          bool
	PressSpace      bool
}

func NewHandler() *Handler {
	return &Handler{
		Protocol:        0,
		HandshakePassed: false,
		PressQ:          false,
		PressW:          false,
		PressSpace:      false,
	}
}

//OnMessage when a message comes in
func (h *Handler) OnMessage(m []byte) (Packet, error) {
	//switch over different packet opcodes
	switch OPCode(m[0]) {
	case OPCodeServerMouseMove:
		p := MouseMove{}
		err := Decode(m, &p)
		return &p, err
	default:
		return nil, fmt.Errorf("unknown packet: %+v", m)
	}
}
