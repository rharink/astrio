package packet

type EjectMass struct {
	packet
}

func (p *EjectMass) Bytes() []byte {
	p.OPCode = OPCodeServerEjectMass
	return Encode(p)
}
