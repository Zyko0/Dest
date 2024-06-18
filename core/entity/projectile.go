package entity

import (
	"image"

	"github.com/Zyko0/Alapae/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

type Projectile struct {
	pos    mgl64.Vec3
	dir    mgl64.Vec3
	radius float64
	speed  float64
	dead   bool
}

func NewProjectile(pos, dir mgl64.Vec3, radius, speed float64) *Projectile {
	return &Projectile{
		pos:    pos,
		dir:    dir,
		radius: radius,
		speed:  speed,
	}
}

func (p *Projectile) Update() {
	p.pos = p.pos.Add(p.dir.Mul(p.speed)) // TODO:
}

func (p *Projectile) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, ctx *graphics.Context) ([]ebiten.Vertex, []uint16) {
	pos := ctx.ProjView.Mul4x1(p.pos.Sub(ctx.CameraPosition).Vec4(1))
	// Behind the camera in screen space
	if pos.Z() < 0 {
		return vx, ix
	}
	camRight := ctx.CameraRight.Mul(p.radius)
	camUp := ctx.CameraUp.Mul(p.radius)
	// Append vertices and indices if the quad is visible
	rect := image.Rect(0, 0, int(p.radius), int(p.radius))
	vx, ix = graphics.AppendBillboardVerticesIndices(
		vx, ix, uint16(*index), rect, p.pos.Sub(ctx.CameraPosition), camRight, camUp, &ctx.ProjView,
	)
	*index++

	return vx, ix
}

func (p *Projectile) Position() mgl64.Vec3 {
	return p.pos
}

func (p *Projectile) Radius() float64 {
	return p.radius
}

func (p *Projectile) Dead() bool {
	return p.dead
}
