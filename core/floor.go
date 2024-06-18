package core

import (
	"github.com/Zyko0/Alapae/assets"
	"github.com/Zyko0/Alapae/core/aoe"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	markerResolutionFactor = 2.
)

type Floor struct {
	markers []*aoe.Marker

	Image *ebiten.Image
}

func newFloor() *Floor {
	return &Floor{
		Image: ebiten.NewImage(
			ArenaSize*markerResolutionFactor,
			ArenaSize*markerResolutionFactor,
		),
	}
}

func (f *Floor) AddMarker(m *aoe.Marker) {
	f.markers = append(f.markers, m)
}

func (f *Floor) Update() {
	var n int
	for _, m := range f.markers {
		m.Update()
		if m.Over() {
			continue
		}
		f.markers[n] = m
		n++
	}
	f.markers = f.markers[:n]
}

func (f *Floor) Draw() {
	// Reset floor texture
	f.Image.Clear()
	// Draw markers
	var vx []ebiten.Vertex
	var ix []uint16
	var index int
	for _, m := range f.markers {
		vx, ix = m.AppendVerticesIndices(vx, ix, &index, markerResolutionFactor)
	}
	f.Image.DrawTrianglesShader(vx, ix, assets.ShaderMarker(), &ebiten.DrawTrianglesShaderOptions{
		Blend: ebiten.Blend{
			BlendFactorSourceRGB:        ebiten.BlendFactorSourceColor,
			BlendFactorSourceAlpha:      ebiten.BlendFactorSourceAlpha,
			BlendFactorDestinationRGB:   ebiten.BlendFactorDestinationColor,
			BlendFactorDestinationAlpha: ebiten.BlendFactorDestinationAlpha,
			BlendOperationRGB:           ebiten.BlendOperationMax,
			BlendOperationAlpha:         ebiten.BlendOperationMax,
		},
	})
}
