package core

import (
	"github.com/Zyko0/Alapae/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

type Entity interface {
	Update()
	AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, ctx *graphics.Context) ([]ebiten.Vertex, []uint16)

	Position() mgl64.Vec3
	Radius() float64
	Dead() bool
}

type Boss interface{}
