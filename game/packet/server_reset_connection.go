package packet

type ResetConnection1 struct {
	packet
	Protocol uint32
}

func (p *ResetConnection1) Bytes() []byte {
	p.OPCode = OPCodeServerResetConnection1
	return Encode(p)
}

type ResetConnection2 struct {
	packet
	Protocol uint32
}

func (p *ResetConnection2) Bytes() []byte {
	p.OPCode = OPCodeServerResetConnection2
	return Encode(p)
}
