package packet

type Split struct {
	packet
}

func (p *Split) Bytes() []byte {
	p.OPCode = OPCodeServerSplit
	return Encode(p)
}
