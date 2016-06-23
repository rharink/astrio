package packet

type SetBorder struct {
	packet
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

func (p *SetBorder) Bytes() []byte {
	p.OPCode = OPCodeClientSetBorder
	return Encode(p)
}
