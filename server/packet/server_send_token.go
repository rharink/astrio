package packet

type SendToken struct {
	packet
	Token [256]byte
}

func (p *SendToken) Bytes() []byte {
	p.OPCode = OPCodeServerSendToken
	return Encode(p)
}
