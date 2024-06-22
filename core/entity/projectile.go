package entity

import (
	"image"
	"image/color"

	"github.com/Zyko0/Alapae/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

const ProjectileMaxDuration = 5 * 60

type Projectile struct {
	pos    mgl64.Vec3
	dir    mgl64.Vec3
	clr0   float32
	clr1   float32
	radius float64
	speed  float64

	ticks uint
	dead  bool
}

func NewProjectile(pos, dir mgl64.Vec3, radius, speed float64, clr0, clr1 color.Color) *Projectile {
	return &Projectile{
		pos:    pos,
		dir:    dir,
		clr0:   graphics.ColorAsFloat32RGB(clr0),
		clr1:   graphics.ColorAsFloat32RGB(clr1),
		radius: radius,
		speed:  speed,
	}
}

func (p *Projectile) Update(_ *Context) {
	p.pos = p.pos.Add(p.dir.Mul(p.speed)) // TODO:
	p.ticks++
}

func (p *Projectile) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, ctx *graphics.Context) ([]ebiten.Vertex, []uint16) {
	pos := ctx.ProjView.Mul4x1(p.pos.Sub(ctx.CameraPosition).Vec4(1))
	// Behind the camera in screen space
	if pos.Z() < 0 {
		return vx, ix
	}
	camRight := ctx.CameraRight.Mul(p.radius * graphics.SpriteScale)
	camUp := ctx.CameraUp.Mul(p.radius * graphics.SpriteScale)
	// Append vertices and indices if the quad is visible
	rect := image.Rect(0, 0, int(p.radius), int(p.radius))
	vi := len(vx)
	vx, ix = graphics.AppendBillboardVerticesIndices(
		vx, ix, uint16(*index), rect, p.pos.Sub(ctx.CameraPosition), camRight, camUp, &ctx.ProjView,
	)
	*index++
	for i := 0; i < 4; i++ {
		vx[vi+i].ColorR = 1 // Bullet hardcoded
		vx[vi+i].ColorG = p.clr0
		vx[vi+i].ColorB = p.clr1
	}

	return vx, ix
}

func (p *Projectile) Position() mgl64.Vec3 {
	return p.pos
}

func (p *Projectile) Radius() float64 {
	return p.radius
}

func (p *Projectile) Dead() bool {
	return p.ticks >= ProjectileMaxDuration || p.pos.Y()-p.radius < -graphics.SpriteScale // TODO: do in collisions check
}
