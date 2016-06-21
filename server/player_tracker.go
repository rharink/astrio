package server

type PlayerTracker struct {
	Mouse struct {
		X float32
		Y float32
	}
	color struct {
		R int
		G int
		B int
	}
	score uint32
	scale int
}
