package graphics

import (
	"image"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	quadVertices = []mgl64.Vec3{
		{-0.5, 0.5, 0},
		{0.5, 0.5, 0},
		{-0.5, -0.5, 0},
		{0.5, -0.5, 0},
	}
)

func AppendBillboardUVVerticesIndices(vx []ebiten.Vertex, ix []uint16, index uint16, rect image.Rectangle, pos, right, up mgl64.Vec3, tf *mgl64.Mat4, truncateY bool) ([]ebiten.Vertex, []uint16) {
	v0 := pos.Add(right.Mul(quadVertices[0].X())).Add(
		up.Mul(quadVertices[0].Y()),
	).Vec4(1)
	v1 := pos.Add(right.Mul(quadVertices[1].X())).Add(
		up.Mul(quadVertices[1].Y()),
	).Vec4(1)
	v2 := pos.Add(right.Mul(quadVertices[2].X())).Add(
		up.Mul(quadVertices[2].Y()),
	).Vec4(1)
	v3 := pos.Add(right.Mul(quadVertices[3].X())).Add(
		up.Mul(quadVertices[3].Y()),
	).Vec4(1)
	if truncateY {
		v0[1] = max(v0[1], -10)
		v1[1] = max(v1[1], -10)
		v2[1] = max(v2[1], -10)
		v3[1] = max(v3[1], -10)
	}
	v0Len := v0.Vec3().Len()
	v1Len := v1.Vec3().Len()
	v2Len := v2.Vec3().Len()
	v3Len := v3.Vec3().Len()
	v0 = tf.Mul4x1(v0)
	v1 = tf.Mul4x1(v1)
	v2 = tf.Mul4x1(v2)
	v3 = tf.Mul4x1(v3)
	for i := 0; i < 2; i++ {
		v0[i] /= v0[3]
		v1[i] /= v1[3]
		v2[i] /= v2[3]
		v3[i] /= v3[3]
	}

	vx = append(vx, []ebiten.Vertex{
		{
			DstX:   float32(v0.X()),
			DstY:   float32(v0.Y()),
			SrcX:   float32(float64(rect.Min.X) / v0.W()),
			SrcY:   float32(float64(rect.Min.Y) / v0.W()),
			ColorR: float32(1 / v0Len),
			ColorG: float32(1 / v0.W()),
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   float32(v1.X()),
			DstY:   float32(v1.Y()),
			SrcX:   float32(float64(rect.Max.X) / v1.W()),
			SrcY:   float32(float64(rect.Min.Y) / v1.W()),
			ColorR: float32(1 / v1Len),
			ColorG: float32(1 / v1.W()),
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   float32(v2.X()),
			DstY:   float32(v2.Y()),
			SrcX:   float32(float64(rect.Min.X) / v2.W()),
			SrcY:   float32(float64(rect.Max.Y) / v2.W()),
			ColorR: float32(1 / v2Len),
			ColorG: float32(1 / v2.W()),
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   float32(v3.X()),
			DstY:   float32(v3.Y()),
			SrcX:   float32(float64(rect.Max.X) / v3.W()),
			SrcY:   float32(float64(rect.Max.Y) / v3.W()),
			ColorR: float32(1 / v3Len),
			ColorG: float32(1 / v3.W()),
			ColorB: 1,
			ColorA: 1,
		},
	}...)

	ix = append(ix, []uint16{
		rectIndices[0] + 4*index,
		rectIndices[1] + 4*index,
		rectIndices[2] + 4*index,
		rectIndices[3] + 4*index,
		rectIndices[4] + 4*index,
		rectIndices[5] + 4*index,
	}...)

	return vx, ix
}

func AppendBillboardVerticesIndices(vx []ebiten.Vertex, ix []uint16, index uint16, pos, right, up mgl64.Vec3, tf *mgl64.Mat4) ([]ebiten.Vertex, []uint16) {
	v0 := pos.Add(right.Mul(quadVertices[0].X())).Add(
		up.Mul(quadVertices[0].Y()),
	).Vec4(1)
	v0 = tf.Mul4x1(v0)
	v1 := pos.Add(right.Mul(quadVertices[1].X())).Add(
		up.Mul(quadVertices[1].Y()),
	).Vec4(1)
	v1 = tf.Mul4x1(v1)
	v2 := pos.Add(right.Mul(quadVertices[2].X())).Add(
		up.Mul(quadVertices[2].Y()),
	).Vec4(1)
	v2 = tf.Mul4x1(v2)
	v3 := pos.Add(right.Mul(quadVertices[3].X())).Add(
		up.Mul(quadVertices[3].Y()),
	).Vec4(1)
	v3 = tf.Mul4x1(v3)
	for i := 0; i < 2; i++ {
		v0[i] /= v0[3]
		v1[i] /= v1[3]
		v2[i] /= v2[3]
		v3[i] /= v3[3]
	}

	vx = append(vx, []ebiten.Vertex{
		{
			DstX:   float32(v0.X()),
			DstY:   float32(v0.Y()),
			SrcX:   -1,
			SrcY:   -1,
			ColorR: 0,
			ColorG: 0,
			ColorB: 1,
			ColorA: float32(v0.Z()),
		},
		{
			DstX:   float32(v1.X()),
			DstY:   float32(v1.Y()),
			SrcX:   1,
			SrcY:   -1,
			ColorR: 0,
			ColorG: 0,
			ColorB: 1,
			ColorA: float32(v1.Z()),
		},
		{
			DstX:   float32(v2.X()),
			DstY:   float32(v2.Y()),
			SrcX:   -1,
			SrcY:   1,
			ColorR: 0,
			ColorG: 0,
			ColorB: 1,
			ColorA: float32(v2.Z()),
		},
		{
			DstX:   float32(v3.X()),
			DstY:   float32(v3.Y()),
			SrcX:   1,
			SrcY:   1,
			ColorR: 0,
			ColorG: 0,
			ColorB: 1,
			ColorA: float32(v3.Z()),
		},
	}...)

	ix = append(ix, []uint16{
		rectIndices[0] + 4*index,
		rectIndices[1] + 4*index,
		rectIndices[2] + 4*index,
		rectIndices[3] + 4*index,
		rectIndices[4] + 4*index,
		rectIndices[5] + 4*index,
	}...)

	return vx, ix
}
