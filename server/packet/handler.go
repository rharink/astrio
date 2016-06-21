package packet

import "fmt"

type (
	Handler struct {
		Protocol   int
		PressQ     bool
		PressW     bool
		PressSpace bool
	}
)

const (
	SET_DIRECTION uint8 = 16
)

func NewHandler() *Handler {
	return &Handler{
		Protocol:   1,
		PressQ:     false,
		PressW:     false,
		PressSpace: false,
	}
}

func (h *Handler) OnMessage(message []byte) {
	switch uint8(message[0]) {
	case SET_DIRECTION:
		fmt.Printf("moving player:%d to x:%d, y:%d \n", uint32(message[3]), uint32(message[1]), uint32(message[2]))
	default:
		return
	}
}
