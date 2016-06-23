package packet

type QKeyReleased struct {
	packet
}

func (p *QKeyReleased) Bytes() []byte {
	p.OPCode = OPCodeServerQKeyReleased
	return Encode(p)
}
