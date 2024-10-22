package entity

import (
	"image/color"

	"github.com/Zyko0/Dest/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

const dotFrequency = 20

type Projectile struct {
	team       Team
	pos        mgl64.Vec3
	dir        mgl64.Vec3
	clrIn      float32
	clrOut     float32
	alpha      float64
	radius     float64
	speed      float64
	dmg        float64
	pull       float64
	duration   uint
	resistance uint
	homing     bool

	lastHit uint
	ticks   uint
	dead    bool
}

func NewProjectile(pos, dir mgl64.Vec3, radius, speed, dmg, pull float64, team Team, clrIn, clrOut color.Color, alpha float64, duration, resistance uint, homing bool) *Projectile {
	return &Projectile{
		team:       team,
		pos:        pos,
		dir:        dir,
		clrIn:      graphics.ColorAsFloat32RGB(clrIn),
		clrOut:     graphics.ColorAsFloat32RGB(clrOut),
		alpha:      alpha,
		radius:     radius,
		speed:      speed,
		dmg:        max(dmg, 0),
		pull:       pull,
		duration:   duration,
		resistance: resistance,
		homing:     homing,
	}
}

func (p *Projectile) Team() Team {
	return p.team
}

func (p *Projectile) Pull() float64 {
	return p.pull
}

func (p *Projectile) Damage() float64 {
	if p.lastHit != 0 && p.ticks-p.lastHit < dotFrequency {
		return 0
	}
	return p.dmg
}

func (p *Projectile) TakeHit(_ float64) {
	if p.resistance > 0 {
		if p.lastHit == 0 || p.ticks-p.lastHit >= dotFrequency {
			p.lastHit = p.ticks
			p.resistance -= 1
		}
	}
}

func (p *Projectile) Update(ctx *Context) {
	if p.homing && ctx.Boss != nil && !ctx.Boss.Dead() {
		p.dir = ctx.Boss.Position().Sub(p.pos).Normalize()
	}
	p.pos = p.pos.Add(p.dir.Mul(p.speed))
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
	vi := len(vx)
	vx, ix = graphics.AppendBillboardVerticesIndices(
		vx, ix, uint16(*index), p.pos.Sub(ctx.CameraPosition), camRight, camUp, &ctx.ProjView,
	)
	*index++
	for i := 0; i < 4; i++ {
		vx[vi+i].ColorR = 2 // Bullet hardcoded
		vx[vi+i].ColorG = p.clrOut
		vx[vi+i].ColorB = p.clrIn
		vx[vi+i].ColorA = float32(p.alpha)
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
	return p.ticks >= p.duration ||
		p.pos.Y()-p.radius < -graphics.SpriteScale ||
		p.resistance == 0
}
