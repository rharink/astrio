package packet

type QKeyPressed struct {
	packet
}

func (p *QKeyPressed) Bytes() []byte {
	p.OPCode = OPCodeServerQKeyPressed
	return Encode(p)
}
