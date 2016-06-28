package packet

type UpdateNodes struct {
	packet
}

func (p *UpdateNodes) Bytes() []byte {
	p.OPCode = OPCodeClientUpdateNodes
	return Encode(p)
}
