package packet

import "testing"

type TestCase struct {
	Packet Packet
	OPCode OPCode
}

var tests = []TestCase{
	{&SetBorder{}, OPCodeClientSetBorder},
	{&UpdateLeaderboardTeam{}, OPCodeClientUpdateLeaderBoardTeam},
	{&UpdateNodes{}, OPCodeClientUpdateNodes},
	{&EjectMass{}, OPCodeServerEjectMass},
	{&MouseMove{}, OPCodeServerMouseMove},
	{&QKeyPressed{}, OPCodeServerQKeyPressed},
	{&QKeyReleased{}, OPCodeServerQKeyReleased},
	{&ResetConnection1{}, OPCodeServerResetConnection1},
	{&ResetConnection2{}, OPCodeServerResetConnection2},
	//{&SendToken{}, OPCodeServerSendToken},
	{&SetNickname{}, OPCodeServerSetNickname},
	{&Spectate{}, OPCodeServerSpectate},
	{&Split{}, OPCodeServerSplit},
}

func TestEncodeDecode(t *testing.T) {
	for _, test := range tests {
		b := test.Packet.Bytes()
		opcode := OPCode(b[0])
		if opcode != test.OPCode {
			t.Fatalf("expected opcode %d got %d", test.OPCode, opcode)
		}
	}
}
