package packet

type MouseMove struct {
	packet
	X  float64
	Y  float64
	ID uint32
}

func (p *MouseMove) Bytes() []byte {
	p.OPCode = OPCodeServerMouseMove
	return Encode(p)
}
