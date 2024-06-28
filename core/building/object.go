package building

import (
	"github.com/Zyko0/Alapae/core/entity"
	"github.com/Zyko0/Alapae/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

type ItemObject struct {
	pos      mgl64.Vec3
	radius   float64
	picked   bool
	targeted bool

	Item *Item
}

func (i *ItemObject) Team() entity.Team {
	return entity.TeamNone
}

func (i *ItemObject) TakeHit(_ float64) {}

func (i *ItemObject) Damage() float64 {
	return 0
}

func (i *ItemObject) Update(ctx *entity.Context) {
	if i.picked {
		return
	}
	dir := ctx.PlayerPosition.Sub(i.pos).Normalize()
	dot := dir.Dot(ctx.PlayerDirection)
	if dot < -0.99 {
		i.targeted = true
	} else {
		i.targeted = false
	}
}

func (o *ItemObject) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, ctx *graphics.Context) ([]ebiten.Vertex, []uint16) {
	pos := ctx.ProjView.Mul4x1(o.pos.Sub(ctx.CameraPosition).Vec4(1))
	// Behind the camera in screen space
	if pos.Z() < 0 {
		return vx, ix
	}
	camRight := ctx.CameraRight.Mul(2 * o.radius * graphics.SpriteScale)
	camUp := ctx.CameraUp.Mul(2 * o.radius * graphics.SpriteScale)
	// Vertices
	vi := len(vx)
	vx, ix = graphics.AppendBillboardUVVerticesIndices(
		vx, ix, uint16(*index), o.Item.def.Rect,
		o.pos.Sub(ctx.CameraPosition),
		camRight,
		camUp,
		&ctx.ProjView,
		true,
	)
	*index++
	for vi < len(vx) {
		vx[vi].ColorR = 1                                                 // Sprite item hardcoded
		vx[vi].ColorB = graphics.AngleOriginAsFloat32(0, o.Item.def.Rect) // Angle
		vx[vi].ColorA = 1
		vi++
	}

	return vx, ix
}

func (i *ItemObject) Position() mgl64.Vec3 {
	return i.pos
}

func (i *ItemObject) Radius() float64 {
	return i.radius
}

func (i *ItemObject) Dead() bool {
	return i.picked
}
