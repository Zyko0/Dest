package ui

import (
	"image"

	"github.com/Zyko0/Dest/assets"
	"github.com/Zyko0/Dest/graphics"
	"github.com/Zyko0/Dest/logic"
	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/opt"
	"github.com/Zyko0/Ebiary/ui/uiex"
	"github.com/hajimehoshi/ebiten/v2"
)

type sliderBar struct {
	init   bool
	rect   image.Rectangle
	cursor float64
}

func (sb *sliderBar) Draw(area *ebiten.Image) {
	if !sb.init {
		sb.rect = area.Bounds()
		sb.init = true
	}
	vx, ix := graphics.AppendRectVerticesIndices(
		nil, nil, 0, &graphics.RectOpts{
			DstX:      float32(area.Bounds().Min.X),
			DstY:      float32(area.Bounds().Min.Y),
			SrcX:      1,
			SrcY:      1,
			DstWidth:  float32(area.Bounds().Dx()) * float32(sb.cursor),
			DstHeight: float32(area.Bounds().Dy()),
			SrcWidth:  1,
			SrcHeight: 1,
			R:         0,
			G:         0.25,
			B:         0.78,
			A:         1,
		},
	)
	area.DrawTriangles(vx, ix, graphics.BrushImage, nil)
}

type slider struct {
	*ui.Block
	offset image.Point
	sb     *sliderBar
	update func(t float64)
}

func newSlider(offset image.Point, update func(t float64), t float64) *slider {
	s := &slider{
		Block: ui.NewBlock().WithOptions(opt.Block.Options(
			opt.Padding(10),
			opt.Border(0, softWhite),
			ui.WithCustomUpdateFunc(func(b *slider, is ui.InputState) {
				b.SetBorderWidth(0)
				if b.sb.init {
					c := is.Cursor().Sub(b.offset)
					if c.In(b.sb.rect) {
						b.SetBorderWidth(2)
						if is.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 0 {
							b.sb.cursor = float64(c.X-b.sb.rect.Min.X) / float64(b.sb.rect.Dx())
							update(b.sb.cursor)
						}
					}
				}
			}),
		)),
		offset: offset,
		sb: &sliderBar{
			cursor: t,
		},
	}
	s.SetContent(s.sb)

	return s
}

type Settings struct {
	layout *ui.Layout

	offset image.Point
	active bool
}

func newSettings() *Settings {
	s := &Settings{
		offset: image.Pt(
			logic.ScreenWidth/2-256,
			logic.ScreenHeight/2-256,
		),
	}
	s.layout = ui.NewLayout(32, 32, image.Rectangle{})
	s.layout.SetDimensions(512, 512)
	s.layout.Grid().WithOptions(opt.Grid.Options(
		opt.RGB(5, 5, 5),
		opt.Rounding(32),
	))
	s.layout.Grid().Add(
		10, 1, 12, 4, uiex.NewLabel("Settings").WithOptions(
			opt.Label.Text(
				opt.Text.Source(assets.FontSource),
				opt.Text.Size(48),
				opt.Text.AlignCenter(),
			),
		),
	)

	options := ui.NewGrid(4, 3).WithOptions(
		opt.Grid.Options(
			opt.RGB(15, 15, 15),
			opt.Rounding(15),
		),
	)
	options.Add(0, 0, 1, 1, uiex.NewLabel("Mouse").WithOptions(
		opt.Label.Text(
			opt.Text.Color(softWhite),
			opt.Text.AlignLeft(),
			opt.Text.Source(assets.FontSource),
			opt.Text.Size(24),
			opt.Text.PaddingLeft(10),
		),
	))
	options.Add(0, 1, 1, 1, uiex.NewLabel("Music").WithOptions(
		opt.Label.Text(
			opt.Text.Color(softWhite),
			opt.Text.AlignLeft(),
			opt.Text.Source(assets.FontSource),
			opt.Text.Size(24),
			opt.Text.PaddingLeft(10),
		),
	))
	options.Add(0, 2, 1, 1, uiex.NewLabel("SFX").WithOptions(
		opt.Label.Text(
			opt.Text.Color(softWhite),
			opt.Text.AlignLeft(),
			opt.Text.Source(assets.FontSource),
			opt.Text.Size(24),
			opt.Text.PaddingLeft(10),
		),
	))

	options.Add(1, 0, 3, 1, newSlider(s.offset, func(t float64) {
		logic.MouseSensitivity = max(t, 0.01)
	}, logic.MouseSensitivity))
	options.Add(1, 1, 3, 1, newSlider(s.offset, func(t float64) {
		logic.MusicVolume = t
		assets.SetMusicVolume(t)
	}, logic.MusicVolume))
	options.Add(1, 2, 3, 1, newSlider(s.offset, func(t float64) {
		logic.SFXVolume = t
		assets.SetSFXVolume(t)
	}, logic.SFXVolume))

	s.layout.Grid().Add(2, 7, 28, 23, options)

	return s
}

func (s *Settings) Update() {
	// Update layout
	s.layout.Update(s.offset, ui.GetInputState())
}

func (s *Settings) Draw(screen *ebiten.Image) {
	s.layout.Draw(screen)
}
