package pattern

import (
	"image/color"

	"github.com/Zyko0/Dest/assets"
	"github.com/Zyko0/Dest/core/aoe"
	"github.com/Zyko0/Dest/core/entity"
	"github.com/Zyko0/Dest/graphics"
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
			1, s.speed, 5, 0,
			entity.TeamEnemy,
			color.RGBA{255, 192, 192, 255},
			color.RGBA{220, 0, 0, 255},
			1, 5*60, 1, false,
		)
		ctx.Entities = append(ctx.Entities, p)
		assets.PlayBossShoot()
		s.over = true
	}
	s.ticks++
}

func (s *Shoot) Over() bool {
	return s.over
}
