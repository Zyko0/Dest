package boss

import (
	"math"

	"github.com/Zyko0/Alapae/core/boss/pattern"
	"github.com/Zyko0/Alapae/core/entity"
)

func NewMoveTo() PatternInstancier {
	return func(ctx *entity.Context) Pattern {
		target := ctx.PlayerPosition
		target[1] = ctx.Boss.Position().Y()
		return pattern.NewMoveTo(
			ctx.Boss.Position(),
			target,
			2.,
			30,
		)
	}
}

func NewMoveToEdge() PatternInstancier {
	return func(ctx *entity.Context) Pattern {
		pos := ctx.PlayerPosition
		dir := pos.Sub(ctx.Boss.Position()).Normalize()
		tx := 0.
		if dir.X() > 0 {
			tx = 192
		}
		tz := 0.
		if dir.Z() > 0 {
			tz = 192
		}
		edge := pos
		if math.Abs(pos.X()-tx) < math.Abs(pos.Z()-tz) {
			edge[0] = tx
			edge[2] = pos.Z() + math.Abs(pos.X()-tx)*dir.Z()
		} else {
			edge[2] = tz
			edge[0] = pos.X() + math.Abs(pos.Z()-tz)*dir.X()
		}
		edge[1] = ctx.Boss.Position().Y()

		return pattern.NewMoveTo(
			ctx.Boss.Position(),
			edge,
			2.,
			30,
		)
	}
}

func NewShoot() PatternInstancier {
	return func(ctx *entity.Context) Pattern {
		target := ctx.PlayerPosition
		target[1] = ctx.Boss.Position().Y()
		return pattern.NewShoot(
			2.,
			30,
		)
	}
}