package packet

type UpdateLeaderboardTeam struct {
	packet
	Amount uint32
	Score  float32
}

func (p *UpdateLeaderboardTeam) Bytes() []byte {
	p.OPCode = OPCodeClientUpdateLeaderBoardTeam
	return Encode(p)
}
