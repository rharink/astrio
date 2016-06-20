package packet

import (
	"bytes"
	"fmt"
)

type (
	Handler struct {
		Protocol   int
		PressQ     bool
		PressW     bool
		PressSpace bool
	}
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
	buffer := bytes.Buffer{}
	buffer.Write(message)
	b, _ := buffer.ReadByte()
	fmt.Println(uint8(message[0]), b)
}
