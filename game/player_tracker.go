package game

import "github.com/rauwekost/astrio/game/packet"

type PlayerTracker struct {
	Mouse struct {
		X float32
		Y float32
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
