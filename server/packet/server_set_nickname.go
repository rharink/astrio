package packet

type SetNickname struct {
	packet
	Nickname [25]byte
}

func (p *SetNickname) Bytes() []byte {
	p.OPCode = OPCodeServerSetNickname
	return Encode(p)
}
