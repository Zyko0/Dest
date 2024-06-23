package boss

import (
	"image"
	"math"
	"math/rand"

	"github.com/Zyko0/Alapae/assets"
	"github.com/Zyko0/Alapae/core/aoe"
	"github.com/Zyko0/Alapae/core/entity"
	"github.com/Zyko0/Alapae/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SmokeMask struct {
	ticks    uint
	position mgl64.Vec3
	stance   entity.Stance
	marker   *aoe.CircleBorder

	smokes  []mgl64.Vec4
	maskDir mgl64.Vec3

	seq        *Sequence
	phase2init bool

	Health    float64
	MaxHealth float64
}

func NewSmokeMask(position mgl64.Vec3) *SmokeMask {
	const (
		particles = 10
	)

	smokes := make([]mgl64.Vec4, particles)
	smokes[0] = mgl64.Vec4{0, 0, 0.25, 1}
	smokes[1] = mgl64.Vec4{0, 0, -0.25, 1}
	smokes[2] = mgl64.Vec4{-1, 0, 1, 1}
	smokes[3] = mgl64.Vec4{0, -1, 1, 1}
	smokes[4] = mgl64.Vec4{1, 0, 1, 1}
	smokes[5] = mgl64.Vec4{0, 1, 1, 1}
	smokes[6] = mgl64.Vec4{-1, 0, -1, 1}
	smokes[7] = mgl64.Vec4{0, -1, -1, 1}
	smokes[8] = mgl64.Vec4{1, 0, -1, 1}
	smokes[9] = mgl64.Vec4{0, 1, -1, 1}

	return &SmokeMask{
		position: position,
		stance:   entity.StanceIdle,
		marker: &aoe.CircleBorder{
			X:      float32(position.X()),
			Y:      float32(position.Z()),
			Radius: BossRadius * graphics.SpriteScale,
		},

		smokes: smokes,
		maskDir: mgl64.Vec3{
			rand.Float64() - 0.5,
			rand.Float64() - 0.5,
			rand.Float64() - 0.5,
		}.Normalize(),

		Health:    0,
		MaxHealth: 0,
	}
}

func (sm *SmokeMask) phase() int {
	if sm.Health <= sm.MaxHealth/2 {
		return 1
	}
	return 0
}

func (sm *SmokeMask) Update(ctx *entity.Context) {
	if !sm.phase2init && sm.phase() == 1 {
		sm.seq = newSequence(
			NewMoveToEdge(),
			NewRandomWalk(),
			NewShoot(), NewShoot(), NewShoot(), NewShoot(), NewShoot(),
			NewRandomWalk(),
			NewMultiPattern(
				NewComet(), NewComet(), NewComet(),
				NewComet(), NewComet(), NewComet(),
			),
			NewRandomWalk(), NewRandomWalk(), NewRandomWalk(),
		)
		sm.phase2init = true
	}
	if sm.seq == nil {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			sm.seq = newSequence(
				NewMoveToEdge(), NewMoveToEdge(), NewMoveToEdge(),
				NewRandomWalk(),
				NewShoot(), NewRandomWalk(),
				NewShoot(), NewRandomWalk(),
				NewShoot(), NewRandomWalk(),
				NewRandomWalk(),
				NewComet(), NewComet(), NewComet(),
				NewRandomWalk(), NewRandomWalk(),
			)

		}
	}
	if sm.seq != nil {
		sm.seq.Update(ctx)
	}
	sm.marker.X = float32(sm.position.X())
	sm.marker.Y = float32(sm.position.Z())
	sm.ticks++
}

func (sm *SmokeMask) SetStance(stance entity.Stance) {
	sm.stance = stance
}

func (sm *SmokeMask) SetPosition(pos mgl64.Vec3) {
	sm.position = pos
}

func (sm *SmokeMask) Image() *ebiten.Image {
	return assets.MaskSheetImage
}

func (sm *SmokeMask) MarkerShape() *aoe.CircleBorder {
	return sm.marker
}

func (sm *SmokeMask) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, ctx *graphics.Context) ([]ebiten.Vertex, []uint16) {
	pos := ctx.ProjView.Mul4x1(sm.position.Sub(ctx.CameraPosition).Vec4(1))
	// Behind the camera in screen space
	if pos.Z() < 0 {
		return vx, ix
	}
	camRight := ctx.CameraRight.Mul(2 * BossRadius * graphics.SpriteScale)
	camUp := ctx.CameraUp.Mul(2 * BossRadius * graphics.SpriteScale)
	// Smoke billboards
	vi := len(vx)
	for i := range sm.smokes {
		t := 1 + math.Abs(float64((sm.ticks+uint(i)*12)%120)/120-0.5)
		smoke := sm.smokes[i].Vec3().Mul(0.5)
		smoke = smoke.Mul(t * graphics.SpriteScale)
		position := sm.position.Add(smoke)
		vx, ix = graphics.AppendBillboardUVVerticesIndices(
			vx, ix, uint16(*index), assets.MaskSmokeSrc,
			position.Sub(ctx.CameraPosition),
			camRight,
			camUp,
			&ctx.ProjView,
			true,
		)
		*index++
	}
	for vi < len(vx) {
		const smokeAlpha = 0.2
		vx[vi].ColorR = 0                                                     // Sprite
		vx[vi].ColorB = graphics.AngleOriginAsFloat32(0, assets.MaskSmokeSrc) // Angle
		vx[vi].ColorA = smokeAlpha                                            // Alpha
		vi++
	}
	// Mask billboard
	var maskRect image.Rectangle
	switch {
	case sm.stance == entity.StanceIdle && sm.phase() == 0:
		maskRect = assets.MaskIdle0Src
	case sm.stance == entity.StanceHostile && sm.phase() == 0:
		maskRect = assets.MaskHostile0Src
	case sm.stance == entity.StanceIdle && sm.phase() == 1:
		maskRect = assets.MaskIdle1Src
	case sm.stance == entity.StanceHostile && sm.phase() == 1:
		maskRect = assets.MaskHostile1Src
	}

	vi = len(vx)
	if sm.ticks%120 == 0 {
		sm.maskDir = mgl64.Vec3{
			rand.Float64() - 0.5,
			rand.Float64() - 0.5,
			rand.Float64() - 0.5,
		}.Normalize()
	}
	t := math.Abs(math.Sin(float64((sm.ticks)%120) / 120 * math.Pi))
	position := sm.position.Add(sm.maskDir.Mul(t * 2.5))
	camRight = camRight.Mul(0.75)
	camUp = camUp.Mul(0.75)
	vx, ix = graphics.AppendBillboardUVVerticesIndices(
		vx, ix, uint16(*index), maskRect,
		position.Sub(ctx.CameraPosition),
		camRight,
		camUp,
		&ctx.ProjView,
		true,
	)
	*index++
	for vi < len(vx) {
		vx[vi].ColorR = 0                                          // Sprite
		vx[vi].ColorB = graphics.AngleOriginAsFloat32(0, maskRect) // Angle
		vx[vi].ColorA = 1                                          // Alpha
		vi++
	}

	return vx, ix
}

func (sm *SmokeMask) Position() mgl64.Vec3 {
	return sm.position
}

func (sm *SmokeMask) Radius() float64 {
	return BossRadius
}

func (sm *SmokeMask) Dead() bool {
	return false
}
