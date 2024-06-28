package hand

import "github.com/go-gl/mathgl/mgl64"

type Side byte

const (
	Right Side = iota
	Left
)

type Weapon byte

const (
	WeaponFinger Weapon = iota
	WeaponPistol
)

type Hand struct {
	Side     Side
	Anim     *AnimationInstance
	Glow     float32
	Rotation mgl64.Vec3
	Fingers  [5]*Finger
}

func New(side Side) *Hand {
	fingers := [5]*Finger{{}, {}, {}, {}, {}}
	*fingers[0] = *Thumb
	*fingers[1] = *Index
	*fingers[2] = *Middle
	*fingers[3] = *Ring
	*fingers[4] = *Pinky

	return &Hand{
		Side:    side,
		Fingers: fingers,
	}
}

func (h *Hand) ShootAnimation(w Weapon) *Animation {
	if w == WeaponFinger {
		return AnimationShootFinger
	}
	return AnimationShootPistol
}

func (h *Hand) AppendData(data []float32) []float32 {
	for _, f := range h.Fingers {
		points := f.ResolvePoints()
		for _, p := range points {
			data = append(data,
				float32(p.X()),
				float32(p.Y()),
				float32(p.Z()),
				0,
			)
		}
	}

	return data
}

func (h *Hand) AsStep() *HandStep {
	hs := &HandStep{
		Rotation: h.Rotation,
	}
	for i := range hs.Fingers {
		hs.Fingers[i] = h.Fingers[i].AsStep()
	}

	return hs
}

func (h *Hand) ShotRightCoeff() float64 {
	switch h.Side {
	case Right:
		return 1
	case Left:
		return -1
	default:
		return 0
	}
}
