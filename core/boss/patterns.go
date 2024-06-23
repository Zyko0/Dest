package boss

import (
	"math"
	"math/rand"

	"github.com/Zyko0/Alapae/core/boss/pattern"
	"github.com/Zyko0/Alapae/core/entity"
	"github.com/Zyko0/Alapae/logic"
	"github.com/go-gl/mathgl/mgl64"
)

// Meta

type multiPattern struct {
	patterns []Pattern
}

func (mp *multiPattern) Update(ctx *entity.Context) {
	for _, p := range mp.patterns {
		p.Update(ctx)
	}
}

func (mp *multiPattern) Over() bool {
	for _, p := range mp.patterns {
		if !p.Over() {
			return false
		}
	}
	return true
}

func NewMultiPattern(instanciers ...PatternInstancier) PatternInstancier {
	return func(ctx *entity.Context) Pattern {
		ps := make([]Pattern, len(instanciers))
		for i := range instanciers {
			ps[i] = instanciers[i](ctx)
		}
		return &multiPattern{
			patterns: ps,
		}
	}
}

// Patterns

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

func NewComet() PatternInstancier {
	return func(ctx *entity.Context) Pattern {
		pos := mgl64.Vec3{
			rand.Float64() * logic.ArenaSize,
			100,
			rand.Float64() * logic.ArenaSize,
		}
		return pattern.NewComet(pos, 1.5, 2, 30)
	}
}

func NewCometTargeted() PatternInstancier {
	return func(ctx *entity.Context) Pattern {
		pos := ctx.PlayerPosition
		pos = mgl64.Vec3{
			pos.X() + (rand.Float64()*logic.ArenaSize - logic.ArenaSize/2),
			100,
			pos.Z() + (rand.Float64()*logic.ArenaSize - logic.ArenaSize/2),
		}
		pos[0] = max(min(pos[0], logic.ArenaSize), 0)
		pos[2] = max(min(pos[2], logic.ArenaSize), 0)
		return pattern.NewComet(pos, 1.5, 2, 30)
	}
}

func NewRandomWalk() PatternInstancier {
	return func(ctx *entity.Context) Pattern {
		target := ctx.PlayerPosition
		target[1] = ctx.Boss.Position().Y()
		return pattern.NewRandomWalk(60, 1)
	}
}
