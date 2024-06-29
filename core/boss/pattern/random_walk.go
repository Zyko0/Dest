package pattern

import (
	"math/rand"

	"github.com/Zyko0/Dest/core/entity"
	"github.com/Zyko0/Dest/logic"
	"github.com/go-gl/mathgl/mgl64"
)

type RandomWalk struct {
	ticks uint
	over  bool

	pos      mgl64.Vec3
	dir      mgl64.Vec3
	length   float64
	traveled float64
	speed    float64
}

func NewRandomWalk(length, speed float64) *RandomWalk {
	return &RandomWalk{
		length: length,
		speed:  speed,
	}
}

func (rw *RandomWalk) Update(ctx *entity.Context) {
	if rw.over {
		return
	}
	if rw.ticks == 0 {
		rw.pos = ctx.Boss.Position()
		x, z := -1., -1.
		for x < 0 || x > logic.ArenaSize || z < 0 || z > logic.ArenaSize {
			rw.dir = mgl64.Vec3{
				rand.Float64() - 0.5, 0, rand.Float64() - 0.5,
			}.Normalize()
			x, z = rw.pos[0]+rw.dir[0]*rw.length, rw.pos[2]+rw.dir[2]*rw.length
		}
	}
	// Move
	ctx.Boss.SetStance(entity.StanceIdle)
	delta := rw.dir
	if rw.traveled+rw.speed > rw.length {
		delta = delta.Mul(rw.length - rw.traveled)
		rw.traveled = rw.length
		rw.over = true
	} else {
		delta = delta.Mul(rw.speed)
		rw.traveled += rw.speed
	}
	rw.pos = rw.pos.Add(delta)
	ctx.Boss.SetPosition(rw.pos)
	rw.ticks++
}

func (rw *RandomWalk) Over() bool {
	return rw.over
}
