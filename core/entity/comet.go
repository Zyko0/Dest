package entity

import (
	"math"
	"math/rand"

	"github.com/Zyko0/Alapae/assets"
	"github.com/Zyko0/Alapae/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

type Comet struct {
	angle    float64
	position mgl64.Vec3
	radius   float64
	speed    float64

	ticks uint
	dead  bool
}

func NewComet(pos mgl64.Vec3, radius, speed float64) *Comet {
	return &Comet{
		angle:    float64(rand.Intn(2) - 1),
		position: pos,
		radius:   radius,
		speed:    speed,
	}
}

func (c *Comet) Team() Team {
	return TeamEnemy
}

func (c *Comet) Damage() float64 {
	return 10
}

func (c *Comet) TakeHit(_ float64) {}

func (c *Comet) Update(_ *Context) {
	c.position = c.position.Add(mgl64.Vec3{0, -1, 0}.Mul(c.speed)) // TODO:
	c.ticks++
}

func (c *Comet) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, ctx *graphics.Context) ([]ebiten.Vertex, []uint16) {
	pos := ctx.ProjView.Mul4x1(c.position.Sub(ctx.CameraPosition).Vec4(1))
	// Behind the camera in screen space
	if pos.Z() < 0 {
		return vx, ix
	}
	camRight := ctx.CameraRight.Mul(2 * c.radius * graphics.SpriteScale)
	camUp := ctx.CameraUp.Mul(2 * c.radius * graphics.SpriteScale)
	// Smokes
	vi := len(vx)
	off := mgl64.Vec3{0.5, 0, 0.5}
	for i := 0; i < 12; i++ {
		vx, ix = graphics.AppendBillboardUVVerticesIndices(
			vx, ix, uint16(*index), assets.MaskCometSmokeSrc,
			c.position.Sub(ctx.CameraPosition).Add(off.Mul(graphics.SpriteScale)),
			camRight.Mul(1.75),
			camUp.Mul(1.75),
			&ctx.ProjView,
			true,
		)
		if i%4 == 0 {
			off[1] += 1.5
		}
		if i%2 == 0 {
			off[0] *= -1
		} else {
			off[2] *= -1
		}
		*index++
	}
	for vi < len(vx) {
		const smokeAlpha = 0.15
		t := math.Abs(float64((c.ticks+uint(vi/4*5))%45)/45*2-1) * 0.05
		vx[vi].ColorR = 0                                                          // Sprite boss hardcoded
		vx[vi].ColorB = graphics.AngleOriginAsFloat32(t, assets.MaskCometSmokeSrc) // Angle
		vx[vi].ColorA = smokeAlpha
		vi++
	}
	// Ball
	vi = len(vx)
	vx, ix = graphics.AppendBillboardUVVerticesIndices(
		vx, ix, uint16(*index), assets.MaskCometBallSrc,
		c.position.Sub(ctx.CameraPosition),
		camRight,
		camUp,
		&ctx.ProjView,
		true,
	)
	*index++
	for i := 0; i < 4; i++ {
		t := math.Abs(c.angle + float64(c.ticks%30)/30)
		vx[vi+i].ColorR = 0                                                         // Sprite boss hardcoded
		vx[vi+i].ColorB = graphics.AngleOriginAsFloat32(t, assets.MaskCometBallSrc) // Angle
		vx[vi+i].ColorA = 1                                                         // Comet alpha
	}

	return vx, ix
}

func (c *Comet) Position() mgl64.Vec3 {
	return c.position
}

func (c *Comet) Radius() float64 {
	return c.radius
}

func (c *Comet) Dead() bool {
	return c.ticks >= 5*60 || c.position.Y() < -c.radius*12*graphics.SpriteScale
}
