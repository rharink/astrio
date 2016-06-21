package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type WorldUpdatePacket struct {
	Opcode uint8
	Eats   uint16
}

func (p *WorldUpdatePacket) Build() io.ReadSeeker {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, p)
	fmt.Println(buf.Bytes())
}

func (p *WorldUpdatePacket) Opcode() uint8 {
	return uint8(0x10)
}
