package pattern

import (
	"image/color"

	"github.com/Zyko0/Alapae/core/aoe"
	"github.com/Zyko0/Alapae/core/entity"
	"github.com/Zyko0/Alapae/graphics"
)

type Shoot struct {
	ticks uint
	over  bool

	speed float64
	delay uint
}

func NewShoot(speed float64, delay uint) *Shoot {
	s := &Shoot{
		speed: speed,
		delay: delay,
	}

	return s
}

func (s *Shoot) Update(ctx *entity.Context) {
	if s.over {
		return
	}
	// Shots marking delay
	if s.ticks == 0 {
		ctx.Markers = append(ctx.Markers, aoe.NewMarker(
			&aoe.Circle{
				X:      float32(ctx.Boss.Position().X()),
				Y:      float32(ctx.Boss.Position().Z()),
				Radius: float32(ctx.Boss.Radius() * graphics.SpriteScale),
			},
			s.delay, 0,
		))
		ctx.Boss.SetStance(entity.StanceHostile)
	}
	if s.ticks >= s.delay {
		dir := ctx.PlayerPosition.Sub(ctx.Boss.Position()).Normalize()
		p := entity.NewProjectile(
			ctx.Boss.Position().Add(dir.Mul(2.5)),
			dir,
			1, s.speed,
			color.RGBA{220, 0, 0, 255},
			color.RGBA{255, 192, 192, 255},
		)
		ctx.Entities = append(ctx.Entities, p)
		s.over = true
	}
	s.ticks++
}

func (s *Shoot) Over() bool {
	return s.over
}
