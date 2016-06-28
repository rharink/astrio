package packet

type Spectate struct {
	packet
}

func (p *Spectate) Bytes() []byte {
	p.OPCode = OPCodeServerSpectate
	return Encode(p)
}
