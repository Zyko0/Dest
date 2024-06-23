package pattern

import (
	"github.com/Zyko0/Alapae/core/aoe"
	"github.com/Zyko0/Alapae/core/entity"
	"github.com/Zyko0/Alapae/graphics"
	"github.com/go-gl/mathgl/mgl64"
)

type Comet struct {
	ticks uint
	over  bool

	pos    mgl64.Vec3
	radius float64
	speed  float64
	delay  uint
}

func NewComet(pos mgl64.Vec3, radius, speed float64, delay uint) *Comet {
	d := &Comet{
		pos:    pos,
		radius: radius,
		speed:  speed,
		delay:  delay,
	}

	return d
}

func (d *Comet) Update(ctx *entity.Context) {
	if d.over {
		return
	}
	// Drop marking delay
	if d.ticks == 0 {
		ctx.Markers = append(ctx.Markers, aoe.NewMarker(
			&aoe.Circle{
				X:      float32(d.pos.X()),
				Y:      float32(d.pos.Z()),
				Radius: float32(d.radius * graphics.SpriteScale),
			},
			90, 0,
		))
		ctx.Boss.SetStance(entity.StanceHostile)
	}
	if d.ticks == d.delay {
		p := entity.NewComet(
			d.pos,
			d.radius,
			d.speed,
		)
		ctx.Entities = append(ctx.Entities, p)
		d.over = true
	}
	d.ticks++
}

func (s *Comet) Over() bool {
	return s.over
}
