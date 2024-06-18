package assets

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	CursorImage *ebiten.Image
)

func init() {
	CursorImage = ebiten.NewImage(16, 16)
	vector.StrokeLine(CursorImage, 0, 0, 16, 16, 4, color.RGBA{0, 255, 0, 255}, true)
	vector.StrokeLine(CursorImage, 0, 16, 16, 0, 4, color.RGBA{0, 255, 0, 255}, true)
	CursorImage.SubImage(image.Rect(3, 3, 13, 13)).(*ebiten.Image).Fill(color.RGBA{})
}
