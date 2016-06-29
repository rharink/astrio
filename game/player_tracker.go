package game

import "github.com/rauwekost/astrio/game/packet"

type PlayerTracker struct {
	Mouse     Position
	Color     Color
	CenterPos Position
	Viewbox   struct {
		MinX       float32
		MinY       float32
		MaxX       float32
		MaxY       float32
		Width      float32
		Height     float32
		HalfWidth  float32
		HalfHeight float32
	}
}

func NewPlayerTracker() *PlayerTracker {
	return &PlayerTracker{}
}

func (t *PlayerTracker) Update(p interface{}) {
	switch p.(type) {
	case *packet.MouseMove:
		t.UpdateMouse(p.(*packet.MouseMove))
	default:
		break
	}
}

func (t *PlayerTracker) UpdateMouse(p *packet.MouseMove) {
	t.Mouse.X = p.X
	t.Mouse.Y = p.Y
}
