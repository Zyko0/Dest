package hand

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Bone struct {
	XY     float64
	Z      float64
	Length float64
}

type Finger struct {
	Base  mgl64.Vec3
	Bones [3]Bone
}

var (
	Thumb = &Finger{
		Base: mgl64.Vec3{-0.8, 0.9, 0},
		Bones: [3]Bone{
			{
				Z:      math.Pi / 2,
				Length: 0.4,
			},
			{
				Z:      0, //math.Pi / 2,
				Length: 0.3,
			},
			{
				Z:      0,
				Length: 0.2,
			},
		},
	}
	Index = &Finger{
		Base: mgl64.Vec3{-0.65, 0.2, 0},
		Bones: [3]Bone{
			{
				Z:      math.Pi / 2,
				Length: 0.5,
			},
			{
				Z:      0,
				Length: 0.4,
			},
			{
				Z:      0,
				Length: 0.3,
			},
		},
	}
	Middle = &Finger{
		Base: mgl64.Vec3{-0.2, 0.1, 0},
		Bones: [3]Bone{
			{
				Z:      math.Pi / 2,
				Length: 0.5,
			},
			{
				Z:      0,
				Length: 0.4,
			},
			{
				Z:      0,
				Length: 0.3,
			},
		},
	}
	Ring = &Finger{
		Base: mgl64.Vec3{0.25, 0.1, 0},
		Bones: [3]Bone{
			{
				Z:      math.Pi / 2,
				Length: 0.5,
			},
			{
				Z:      0,
				Length: 0.4,
			},
			{
				Z:      0,
				Length: 0.3,
			},
		},
	}
	Pinky = &Finger{
		Base: mgl64.Vec3{0.65, 0.4, 0},
		Bones: [3]Bone{
			{
				Z:      math.Pi / 2,
				Length: 0.4,
			},
			{
				Z:      0,
				Length: 0.3,
			},
			{
				Z:      0,
				Length: 0.3,
			},
		},
	}
)

func (f *Finger) ResolvePoints() [4]mgl64.Vec3 {
	points := [4]mgl64.Vec3{f.Base}
	var xy, z float64
	for i, b := range f.Bones {
		xy += b.XY
		z += b.Z
		x, y := math.Sincos(xy)
		yy, zz := math.Sincos(z)
		dir := mgl64.Vec3{x, y * yy, zz}.Normalize()
		points[i+1] = points[i].Add(dir.Mul(-b.Length))
	}

	return points
}

func (f *Finger) AsStep() *FingerStep {
	return &FingerStep{
		{
			XY: f.Bones[0].XY,
			Z:  f.Bones[0].Z,
		},
		{
			XY: f.Bones[1].XY,
			Z:  f.Bones[1].Z,
		},
		{
			XY: f.Bones[2].XY,
			Z:  f.Bones[2].Z,
		},
	}
}
