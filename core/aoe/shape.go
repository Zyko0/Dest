package aoe

import (
	"math"

	"github.com/Zyko0/Alapae/graphics"
	"github.com/hajimehoshi/ebiten/v2"
)

type Shape interface {
	Update()

	AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, factor, intensity float32) ([]ebiten.Vertex, []uint16)
}

type ShapeType byte

const (
	ShapeFilled ShapeType = iota
	ShapeCircle
	ShapeArrow
	ShapeXCross
	ShapeCircleBorder
)

// Circle

type Circle struct {
	X      float32
	Y      float32
	Radius float32
}

func (c *Circle) Update() {}

func (c *Circle) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, factor, intensity float32) ([]ebiten.Vertex, []uint16) {
	x, y, r := c.X*factor, c.Y*factor, c.Radius*factor
	vx, ix = graphics.AppendRectVerticesIndices(vx, ix, *index, &graphics.RectOpts{
		DstX:      x - r,
		DstY:      y - r,
		SrcX:      -1,
		SrcY:      -1,
		DstWidth:  r * 2,
		DstHeight: r * 2,
		SrcWidth:  2,
		SrcHeight: 2,
		R:         float32(ShapeCircle),
		G:         intensity,
		B:         0,
		A:         0,
	})
	*index++
	return vx, ix
}

// Arrow

type Arrow struct {
	X        float32
	Y        float32
	Size     float32
	Rotation float32
}

func (a *Arrow) Update() {}

func (a *Arrow) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, factor, intensity float32) ([]ebiten.Vertex, []uint16) {
	x, y, r := a.X*factor, a.Y*factor, a.Size/2*factor
	vx, ix = graphics.AppendRectVerticesIndices(vx, ix, *index, &graphics.RectOpts{
		DstX:      x - r,
		DstY:      y - r,
		SrcX:      -1,
		SrcY:      -1,
		DstWidth:  r * 2,
		DstHeight: r * 2,
		SrcWidth:  2,
		SrcHeight: 2,
		R:         float32(ShapeArrow),
		G:         intensity,
		B:         a.Rotation,
		A:         0,
	})
	*index++

	return vx, ix
}

// XCross

type XCross struct {
	Size         float32
	Radius       float32
	Rotation     float32
	RotationIncr float32
}

func (c *XCross) Update() {
	c.Rotation += c.RotationIncr
	if c.Rotation > math.Pi {
		c.Rotation -= math.Pi
	}
}

func (c *XCross) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, factor, intensity float32) ([]ebiten.Vertex, []uint16) {
	vx, ix = graphics.AppendRectVerticesIndices(vx, ix, *index, &graphics.RectOpts{
		DstX:      0,
		DstY:      0,
		SrcX:      -1,
		SrcY:      -1,
		DstWidth:  c.Size * factor,
		DstHeight: c.Size * factor,
		SrcWidth:  2,
		SrcHeight: 2,
		R:         float32(ShapeXCross),
		G:         intensity,
		B:         c.Rotation,
		A:         c.Radius,
	})
	*index++

	return vx, ix
}

// Circle border

type CircleBorder Circle

func (cb *CircleBorder) Update() {}

func (cb *CircleBorder) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, factor, intensity float32) ([]ebiten.Vertex, []uint16) {
	x, y, r := cb.X*factor, cb.Y*factor, cb.Radius*factor
	vx, ix = graphics.AppendRectVerticesIndices(vx, ix, *index, &graphics.RectOpts{
		DstX:      x - r,
		DstY:      y - r,
		SrcX:      -1,
		SrcY:      -1,
		DstWidth:  r * 2,
		DstHeight: r * 2,
		SrcWidth:  2,
		SrcHeight: 2,
		R:         float32(ShapeCircleBorder),
		G:         intensity,
		B:         0,
		A:         0,
	})
	*index++
	
	return vx, ix
}
