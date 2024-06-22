package assets

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	CursorImage *ebiten.Image

	//go:embed images
	fsys embed.FS

	MaskSheetImage    *ebiten.Image
	MaskSmokeSrc      = image.Rect(0, 0, 512, 512)
	MaskIdle0Src      = image.Rect(512, 0, 768, 256)
	MaskHostile0Src   = image.Rect(768, 0, 1024, 256)
	MaskIdle1Src      = image.Rect(512, 256, 768, 512)
	MaskHostile1Src   = image.Rect(768, 256, 1024, 512)
	MaskCometBallSrc  = image.Rect(1024, 0, 1280, 256)
	MaskCometSmokeSrc = image.Rect(1024, 256, 1280, 512)
)

func init() {
	CursorImage = ebiten.NewImage(16, 16)
	vector.StrokeLine(CursorImage, 0, 0, 16, 16, 4, color.RGBA{0, 255, 0, 255}, true)
	vector.StrokeLine(CursorImage, 0, 16, 16, 0, 4, color.RGBA{0, 255, 0, 255}, true)
	CursorImage.SubImage(image.Rect(3, 3, 13, 13)).(*ebiten.Image).Fill(color.RGBA{})

	// Sprites

	b, _ := fsys.ReadFile("images/mask_sheet.png")
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal("err: ", err)
	}
	MaskSheetImage = ebiten.NewImageFromImage(img)
}
