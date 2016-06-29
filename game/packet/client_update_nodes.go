package packet

type UpdateNodes struct {
	packet
	Nodes []Node
}

type Node struct {
	ID    uint16
	X     float32
	Y     float32
	Size  uint16
	R     uint8
	G     uint8
	B     uint8
	Flags uint8
	Name  string
}

func (p *UpdateNodes) Bytes() []byte {
	p.OPCode = OPCodeClientUpdateNodes
	w := NewBinaryWriter()

	w.WriteUint8(uint8(OPCodeClientUpdateNodes))
	for _, n := range p.Nodes {
		w.WriteUint16(n.ID)
		w.WriteFloat(n.X)
		w.WriteFloat(n.Y)
		w.WriteUint16(n.Size)
		w.WriteUint8(n.R)
		w.WriteUint8(n.G)
		w.WriteUint8(n.B)
		w.WriteUint8(n.Flags)
		w.WriteBytes([]byte(n.Name))
		w.WriteUint16(0)
	}

	return w.Bytes()
}
