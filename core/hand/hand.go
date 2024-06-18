package hand

import "github.com/go-gl/mathgl/mgl64"

type HandIndices struct {
	Thumb  int
	Index  int
	Middle int
	Ring   int
	Pinky  int
}

var (
	Right = &HandIndices{
		Thumb:  0,
		Index:  1,
		Middle: 2,
		Ring:   3,
		Pinky:  4,
	}
)

type Hand struct {
	Rotation mgl64.Vec3
	Fingers  [5]*Finger
}

func New() *Hand {
	fingers := [5]*Finger{{}, {}, {}, {}, {}}
	*fingers[0] = *Thumb
	*fingers[1] = *Index
	*fingers[2] = *Middle
	*fingers[3] = *Ring
	*fingers[4] = *Pinky

	return &Hand{
		Fingers: fingers,
	}
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
