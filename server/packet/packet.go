package packet

import "io"

type Packet interface {
	Build() io.ReadSeeker
	Opcode() uint8
}
