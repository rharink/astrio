package game

type CellType uint8

const (
	CellTypePlayer      CellType = 0
	CellTypeFood        CellType = 1
	CellTypeVirus       CellType = 2
	CellTypeEjectedMass CellType = 3
)

//Cell ...
type Cell struct {
	Player         *Player
	TOB            uint16 //tick of birth, to determine age
	Color          Color
	Position       Position
	Size           uint16
	Mass           uint16
	Type           CellType
	KilledBy       *Cell
	IsMoving       bool
	BoostDistance  uint16
	BoostDirection struct {
		X     float32
		Y     float32
		Angle float32
	}
}

//NewCell returns a new cell
func NewCell(typ CellType, tick uint16, p *Player) *Cell {
	c := Cell{
		Player: p,
		TOB:    tick,
		Type:   typ,
	}

	return &c
}
