package entity

import (
	"image"

	"github.com/Zyko0/Alapae/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

const ImpactMaxDuration = 60

type Impact struct {
	pos    mgl64.Vec3
	dir    mgl64.Vec3
	radius float64
	speed  float64

	ticks uint
	dead  bool
}

func NewImpact(pos, dir mgl64.Vec3, radius, speed float64) *Impact {
	return &Impact{
		pos:    pos,
		dir:    dir,
		radius: radius,
		speed:  speed,
	}
}

func (i *Impact) Update(_ *Context) {
	i.pos = i.pos.Add(i.dir.Mul(i.speed)) // TODO:
	i.ticks++
}

func (i *Impact) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, ctx *graphics.Context) ([]ebiten.Vertex, []uint16) {
	pos := ctx.ProjView.Mul4x1(i.pos.Sub(ctx.CameraPosition).Vec4(1))
	// Behind the camera in screen space
	if pos.Z() < 0 {
		return vx, ix
	}
	camRight := ctx.CameraRight.Mul(i.radius * graphics.SpriteScale)
	camUp := ctx.CameraUp.Mul(i.radius * graphics.SpriteScale)
	// Append vertices and indices if the quad is visible
	rect := image.Rect(0, 0, int(i.radius), int(i.radius))
	vi := len(vx)
	vx, ix = graphics.AppendBillboardVerticesIndices(
		vx, ix, uint16(*index), rect, i.pos.Sub(ctx.CameraPosition), camRight, camUp, &ctx.ProjView,
	)
	*index++
	for i := 0; i < 4; i++ {
		vx[vi+i].ColorR = 1 // Bullet hardcoded
	}

	return vx, ix
}

func (i *Impact) Position() mgl64.Vec3 {
	return i.pos
}

func (i *Impact) Radius() float64 {
	return i.radius
}

func (i *Impact) Dead() bool {
	return i.ticks >= ProjectileMaxDuration || i.pos.Y()-i.radius < -graphics.SpriteScale // TODO: do in collisions check
}
