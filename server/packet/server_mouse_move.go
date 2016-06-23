package packet

type MouseMove struct {
	packet
	X  float32
	Y  float32
	ID uint32
}

func (p *MouseMove) Bytes() []byte {
	p.OPCode = OPCodeServerMouseMove
	return Encode(p)
}
