package entity

import (
	"github.com/Zyko0/Dest/core/aoe"
	"github.com/Zyko0/Dest/graphics"
	"github.com/Zyko0/Dest/logic"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

type Portal struct {
	ticks uint
	pos   mgl64.Vec3

	active   bool
	targeted bool
}

func NewPortal() *Portal {
	return &Portal{
		pos: mgl64.Vec3{
			logic.ArenaSize / 2,
			3,
			logic.ArenaSize - 4,
		},
		active: false,
	}
}

func (p *Portal) MarkerShape() aoe.Shape {
	if !p.active || !p.targeted {
		return nil
	}

	return &aoe.Circle{
		X:      float32(p.pos.X()),
		Y:      float32(p.pos.Z()),
		Radius: float32(1.5 * graphics.SpriteScale),
	}
}

func (p *Portal) Team() Team {
	return TeamNone
}

func (p *Portal) Damage() float64 {
	return 0
}

func (p *Portal) TakeHit(_ float64) {}

func (p *Portal) Targeted() bool {
	return p.targeted
}

func (p *Portal) Activate() {
	p.active = true
	p.ticks = 0
}

func (p *Portal) Deactivate() {
	p.active = false
}

func (p *Portal) Update(ctx *Context) {
	p.targeted = false
	if !p.active {
		return
	}

	const targetRange = 64 * 64

	dir := ctx.PlayerPosition.Sub(p.pos).Normalize()
	dot := dir.Dot(ctx.PlayerDirection)
	if dot < -0.95 && ctx.PlayerPosition.Sub(p.pos).LenSqr() < targetRange {
		p.targeted = true
	}
	p.ticks++
}

func (p *Portal) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, ctx *graphics.Context) ([]ebiten.Vertex, []uint16) {
	pos := ctx.ProjView.Mul4x1(p.pos.Sub(ctx.CameraPosition).Vec4(1))
	// Behind the camera in screen space
	if pos.Z() < 0 {
		return vx, ix
	}
	camRight := ctx.CameraRight.Mul(1.5 * graphics.SpriteScale)
	camUp := ctx.CameraUp.Mul(3 * graphics.SpriteScale)
	// Vertices
	vi := len(vx)
	vx, ix = graphics.AppendBillboardVerticesIndices(
		vx, ix, uint16(*index),
		p.pos.Sub(ctx.CameraPosition),
		camRight,
		camUp,
		&ctx.ProjView,
	)
	*index++
	for vi < len(vx) {
		vx[vi].ColorR = 3 // Portal item hardcoded
		vx[vi].ColorB = float32(p.ticks) / 60
		vx[vi].ColorA = 1
		vi++
	}

	return vx, ix
}

func (p *Portal) Position() mgl64.Vec3 {
	return p.pos
}

func (p *Portal) Radius() float64 {
	return 0
}

func (p *Portal) Dead() bool {
	return !p.active
}
