package ui

import (
	"github.com/Zyko0/Dest/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	splashDisplayDuration = 60 * 4
)

type SplashView struct {
	ticks  uint
	active bool

	colorsc ebiten.ColorScale
}

func NewSplashView() *SplashView {
	return &SplashView{
		ticks:  0,
		active: true,

		colorsc: ebiten.ColorScale{},
	}
}

func (sv *SplashView) Active() bool {
	return sv.active
}

func (sv *SplashView) Update() {
	sv.ticks++

	if sv.ticks > splashDisplayDuration {
		sv.active = false
		return
	}
	if len(inpututil.AppendPressedKeys(nil)) > 0 || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		sv.active = false
	}

	d := float64(sv.ticks) / float64(splashDisplayDuration)
	sc := float32(-(d*d)+d) * 4
	sv.colorsc.Reset()
	sv.colorsc.Scale(sc, sc, sc, 1.)
}

func (sv *SplashView) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.SplashImage, &ebiten.DrawImageOptions{
		Filter:     ebiten.FilterLinear,
		ColorScale: sv.colorsc,
	})
}
