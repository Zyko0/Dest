package pattern

import (
	"math"

	"github.com/Zyko0/Dest/assets"
	"github.com/Zyko0/Dest/core/aoe"
	"github.com/Zyko0/Dest/core/entity"
	"github.com/Zyko0/Dest/graphics"
	"github.com/go-gl/mathgl/mgl64"
)

type MoveTo struct {
	ticks uint
	over  bool

	pos      mgl64.Vec3
	angle    float64
	target   mgl64.Vec3
	dir      mgl64.Vec3
	length   float64
	traveled float64
	speed    float64
	delay    uint
	charge   bool
}

func NewMoveTo(pos, target mgl64.Vec3, speed float64, delay uint, charge bool) *MoveTo {
	mt := &MoveTo{
		pos:      pos,
		target:   target,
		dir:      target.Sub(pos).Normalize(),
		length:   target.Sub(pos).Len(),
		traveled: 0,
		speed:    speed,
		delay:    delay,
		charge:   charge,
	}
	mt.angle = math.Atan2(mt.dir.X(), mt.dir.Z())

	return mt
}

func (mt *MoveTo) Update(ctx *entity.Context) {
	const (
		arrowTicks  = 30
		arrowExtra  = 30
		arrow0Start = 0
		arrow1Start = 5
		arrow2Start = 10
	)
	if mt.over {
		return
	}
	// Dash marking delay
	if mt.ticks <= mt.delay {
		d := 2.5
		switch mt.ticks {
		case arrow0Start:
			pos := mt.pos.Add(mt.dir.Mul(d * graphics.SpriteScale))
			ctx.Markers = append(ctx.Markers, aoe.NewMarker(
				&aoe.Arrow{
					X:        float32(pos.X()),
					Y:        float32(pos.Z()),
					Size:     float32(2 * ctx.Boss.Radius() * graphics.SpriteScale),
					Rotation: float32(mt.angle),
				},
				arrowTicks, arrowExtra+arrow1Start+arrow2Start,
			))
		case arrow1Start:
			d *= 2
			pos := mt.pos.Add(mt.dir.Mul(d * graphics.SpriteScale))
			ctx.Markers = append(ctx.Markers, aoe.NewMarker(
				&aoe.Arrow{
					X:        float32(pos.X()),
					Y:        float32(pos.Z()),
					Size:     float32(2 * ctx.Boss.Radius() * graphics.SpriteScale),
					Rotation: float32(mt.angle),
				},
				arrowTicks, arrowExtra+arrow2Start,
			))
		case arrow2Start:
			d *= 3
			pos := mt.pos.Add(mt.dir.Mul(d * graphics.SpriteScale))
			ctx.Markers = append(ctx.Markers, aoe.NewMarker(
				&aoe.Arrow{
					X:        float32(pos.X()),
					Y:        float32(pos.Z()),
					Size:     float32(2 * ctx.Boss.Radius() * graphics.SpriteScale),
					Rotation: float32(mt.angle),
				},
				arrowTicks, arrowExtra,
			))
		}
		ctx.Boss.SetStance(entity.StanceIdle)
	}
	// Dash
	if mt.ticks == mt.delay && mt.charge {
		assets.PlayBossCharge()
	}
	if mt.ticks > mt.delay {
		delta := mt.dir
		if mt.traveled+mt.speed > mt.length {
			delta = delta.Mul(mt.length - mt.traveled)
			mt.traveled = mt.length
			mt.over = true
		} else {
			delta = delta.Mul(mt.speed)
			mt.traveled += mt.speed
		}
		mt.pos = mt.pos.Add(delta)
		ctx.Boss.SetPosition(mt.pos)
		ctx.Boss.SetStance(entity.StanceHostile)
	}
	mt.ticks++
}

func (mt *MoveTo) Over() bool {
	return mt.over
}
