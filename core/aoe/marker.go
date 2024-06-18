package aoe

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Marker struct {
	shape    Shape
	ticks    uint
	duration uint
	extra    uint
}

func NewMarker(shape Shape, duration, extra uint) *Marker {
	return &Marker{
		shape:    shape,
		ticks:    0,
		duration: duration,
		extra:    extra,
	}
}

func (m *Marker) Update() {
	m.ticks = min(m.ticks+1, m.duration+m.extra)
	m.shape.Update()
}

func (m *Marker) intensity() float32 {
	return 0.5 + min(float32(m.ticks)/float32(m.duration), 1)*0.5
}

func (m *Marker) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, factor float32) ([]ebiten.Vertex, []uint16) {
	return m.shape.AppendVerticesIndices(vx, ix, index, factor, m.intensity())
}

func (m *Marker) Over() bool {
	return m.ticks >= m.duration+m.extra
}
