package packet

import (
	"bytes"
	"encoding/binary"
)

type OPCode uint8

const (
	//sending these
	OPCodeClientUpdateNodes           OPCode = 16
	OPCodeClientSpectatePositionSize  OPCode = 17
	OPCodeClientClearAll              OPCode = 20
	OPCodeClientDrawLine              OPCode = 21
	OPCodeClientAddNode               OPCode = 32
	OPCodeClientUpdateLeaderBoardFFA  OPCode = 49
	OPCodeClientUpdateLeaderBoardTeam OPCode = 50
	OPCodeClientSetBorder             OPCode = 64

	//receiving these
	OPCodeServerSetNickname      OPCode = 0
	OPCodeServerSpectate         OPCode = 1
	OPCodeServerMouseMove        OPCode = 16
	OPCodeServerSplit            OPCode = 17
	OPCodeServerQKeyPressed      OPCode = 18
	OPCodeServerQKeyReleased     OPCode = 19
	OPCodeServerEjectMass        OPCode = 21
	OPCodeServerSendToken        OPCode = 80
	OPCodeServerResetConnection1 OPCode = 254
	OPCodeServerResetConnection2 OPCode = 255
)

//Packet packet interface. Bytes() creates a byte representation of the packet
type Packet interface {
	Bytes() []byte
}

//packet internal struct that represents a base packet
type packet struct {
	OPCode OPCode
}

//Encode a packet as byte slice
func Encode(v interface{}) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, v)

	return buf.Bytes()
}

//decode a byteslice to a struct
func Decode(b []byte, v interface{}) error {
	buf := new(bytes.Buffer)
	buf.Write(b)

	return binary.Read(buf, binary.LittleEndian, v)
}
