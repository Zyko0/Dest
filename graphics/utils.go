package graphics

import (
	"image"
	"image/color"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	BrushImage  = ebiten.NewImage(3, 3)
	rectIndices = [6]uint16{0, 1, 2, 1, 2, 3}
)

func init() {
	BrushImage.Fill(color.White)
}

type RectOpts struct {
	DstX, DstY          float32
	SrcX, SrcY          float32
	DstWidth, DstHeight float32
	SrcWidth, SrcHeight float32
	R, G, B, A          float32
}

func AppendRectVerticesIndices(vertices []ebiten.Vertex, indices []uint16, index int, opts *RectOpts) ([]ebiten.Vertex, []uint16) {
	sx, sy, dx, dy := opts.SrcX, opts.SrcY, opts.DstX, opts.DstY
	sw, sh, dw, dh := opts.SrcWidth, opts.SrcHeight, opts.DstWidth, opts.DstHeight
	r, g, b, a := opts.R, opts.G, opts.B, opts.A
	vertices = append(vertices, []ebiten.Vertex{
		{
			DstX:   dx,
			DstY:   dy,
			SrcX:   sx,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   dx + dw,
			DstY:   dy,
			SrcX:   sx + sw,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   dx,
			DstY:   dy + dh,
			SrcX:   sx,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   dx + dw,
			DstY:   dy + dh,
			SrcX:   sx + sw,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}...)

	indiceCursor := uint16(index * 4)
	indices = append(indices, []uint16{
		rectIndices[0] + indiceCursor,
		rectIndices[1] + indiceCursor,
		rectIndices[2] + indiceCursor,
		rectIndices[3] + indiceCursor,
		rectIndices[4] + indiceCursor,
		rectIndices[5] + indiceCursor,
	}...)

	return vertices, indices
}

type QuadOpts struct {
	P0, P1, P2, P3 mgl64.Vec4
	R, G, B, A     float32
}

func AppendQuadVerticesIndices(vertices []ebiten.Vertex, indices []uint16, index int, opts *QuadOpts) ([]ebiten.Vertex, []uint16) {
	r, g, b, a := opts.R, opts.G, opts.B, opts.A
	vertices = append(vertices, []ebiten.Vertex{
		{
			DstX:   float32(opts.P0[0]),
			DstY:   float32(opts.P0[1]),
			SrcX:   float32(opts.P0[2]),
			SrcY:   float32(opts.P0[3]),
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   float32(opts.P1[0]),
			DstY:   float32(opts.P1[1]),
			SrcX:   float32(opts.P1[2]),
			SrcY:   float32(opts.P1[3]),
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   float32(opts.P2[0]),
			DstY:   float32(opts.P2[1]),
			SrcX:   float32(opts.P2[2]),
			SrcY:   float32(opts.P2[3]),
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   float32(opts.P3[0]),
			DstY:   float32(opts.P3[1]),
			SrcX:   float32(opts.P3[2]),
			SrcY:   float32(opts.P3[3]),
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}...)

	indiceCursor := uint16(index * 4)
	indices = append(indices, []uint16{
		rectIndices[0] + indiceCursor,
		rectIndices[1] + indiceCursor,
		rectIndices[2] + indiceCursor,
		rectIndices[3] + indiceCursor,
		rectIndices[4] + indiceCursor,
		rectIndices[5] + indiceCursor,
	}...)

	return vertices, indices
}

func ScreenVertices(vertices []ebiten.Vertex, width, height int) {
	for i := range vertices {
		vertices[i].DstX = float32(float64(1.0-vertices[i].DstX) * float64(width-1) / 2.0)
		vertices[i].DstY = float32(float64(1.0-vertices[i].DstY) * float64(height-1) / 2.0)
	}
}

func ScreenCoordinates(x, y float64, rect image.Rectangle) (float32, float32) {
	sx := float32(float64(1.0-x) * float64(rect.Dx()-1) / 2.0)
	sy := float32(float64(1.0-y) * float64(rect.Dy()-1) / 2.0)
	return sx, sy
}

func ColorAsFloat32RGB(clr color.Color) float32 {
	if clr == nil {
		return 0
	}
	r, g, b, _ := clr.RGBA()
	return float32((r&255)<<16 + (g&255)<<8 + b&255)
}

func AngleOriginAsFloat32(angle float64, rect image.Rectangle) float32 {
	const unit = 256
	a := uint32(angle * 255)
	vx := uint32(rect.Min.X / unit)
	vy := uint32(rect.Min.Y / unit)
	dx := uint32(rect.Dx() / unit)
	dy := uint32(rect.Dy() / unit)
	/*fmt.Println("rect:", rect, "vx", vx, "vy", vy, "dx", dx, "dy", dy)
	fmt.Printf("%32b\n", (a&255)<<16+(vx<<12+vy<<8)+(dx<<4+dy))*/
	return float32((a&255)<<16 + (vx<<12 + vy<<8) + (dx<<4 + dy))
}
