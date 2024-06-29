package ui

import (
	"fmt"
	"strings"

	"github.com/Zyko0/Alapae/assets"
	"github.com/Zyko0/Alapae/core"
	"github.com/Zyko0/Alapae/core/building"
	"github.com/Zyko0/Alapae/graphics"
	"github.com/Zyko0/Alapae/logic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type HUD struct {
}

var (
	titleFace = &text.GoTextFace{
		Source: assets.FontSource,
		Size:   48,
	}
	textFace = &text.GoTextFace{
		Source: assets.FontSource,
		Size:   24,
	}
	tinyFace = &text.GoTextFace{
		Source: assets.FontSource,
		Size:   16,
	}
)

type HUDContext struct {
	Stage       int
	StageKind   core.Stage
	PlayerHP    float64
	PlayerMaxHP float64
	BossHP      float64
	BossMaxHP   float64
	TargetItem  *building.Item
}

func (*HUD) drawItemTooltip(screen *ebiten.Image, item *building.Item) {
	const (
		itemCardWidth    = 768
		itemCardHeight   = 256
		cursesCardWidth  = 512
		cursesCardHeight = 832
	)

	// Item card
	vx, ix := graphics.AppendRectVerticesIndices(
		nil, nil, 0, &graphics.RectOpts{
			DstX:      logic.ScreenWidth/2 - itemCardWidth/2,
			DstY:      64,
			SrcX:      1,
			SrcY:      1,
			DstWidth:  itemCardWidth,
			DstHeight: itemCardHeight,
			R:         0.2,
			G:         0.2,
			B:         0.2,
			A:         0.5,
		},
	)
	// Curses
	vx, ix = graphics.AppendRectVerticesIndices(
		vx, ix, 1, &graphics.RectOpts{
			DstX:      logic.ScreenWidth/2 + itemCardWidth/2 + 32,
			DstY:      64,
			SrcX:      1,
			SrcY:      1,
			DstWidth:  cursesCardWidth,
			DstHeight: cursesCardHeight,
			R:         0.4,
			G:         0.2,
			B:         0.2,
			A:         0.5,
		},
	)
	screen.DrawTriangles(vx, ix, graphics.BrushImage, nil)

	// Icons
	const (
		itemPictureSize  = 192
		cursePictureSize = 64
	)
	// Item
	vx, ix = graphics.AppendRectVerticesIndices(
		vx[:0], ix[:0], 0, &graphics.RectOpts{
			DstX:      logic.ScreenWidth/2 - itemCardWidth/2 + 32,
			DstY:      64 + itemCardHeight/2 - itemPictureSize/2,
			SrcX:      float32(item.SourceRect().Min.X),
			SrcY:      float32(item.SourceRect().Min.Y),
			DstWidth:  itemPictureSize,
			DstHeight: itemPictureSize,
			SrcWidth:  float32(item.SourceRect().Dx()),
			SrcHeight: float32(item.SourceRect().Dy()),
			R:         1,
			G:         1,
			B:         1,
			A:         1,
		},
	)
	// Curses
	const maxCurses = 8
	for i, c := range item.Curses[:min(len(item.Curses), maxCurses)] {
		vx, ix = graphics.AppendRectVerticesIndices(
			vx, ix, i+1, &graphics.RectOpts{
				DstX:      logic.ScreenWidth/2 + itemCardWidth/2 + 48,
				DstY:      64 + 16 + float32(i*96),
				SrcX:      float32(c.SourceRect().Min.X),
				SrcY:      float32(c.SourceRect().Min.Y),
				DstWidth:  cursePictureSize,
				DstHeight: cursePictureSize,
				SrcWidth:  float32(c.SourceRect().Dx()),
				SrcHeight: float32(c.SourceRect().Dy()),
				R:         1,
				G:         1,
				B:         1,
				A:         1,
			},
		)
	}
	screen.DrawTriangles(vx, ix, assets.ItemSheetImage, nil)

	// Text
	topts := &text.DrawOptions{}
	// Name
	str := item.Name()
	topts.GeoM.Reset()
	topts.GeoM.Translate(
		logic.ScreenWidth/2-itemCardWidth/2+64+itemPictureSize,
		96,
	)
	topts.ColorScale.Scale(0.8, 0.8, 0.8, 1)
	text.Draw(screen, str, titleFace, topts)
	// Description
	str = ""
	for i, r := range item.Description() {
		str += string(r)
		if i%40 == 0 {
			if idx := strings.LastIndex(str, " "); idx != -1 {
				str = str[:idx] + "\n" + str[idx+1:]
			}
		}
	}
	topts.LineSpacing = 1.5 * textFace.Size
	topts.GeoM.Reset()
	topts.GeoM.Translate(
		logic.ScreenWidth/2-itemCardWidth/2+64+itemPictureSize,
		156,
	)
	text.Draw(screen, str, textFace, topts)
	// Curses
	for i, c := range item.Curses {
		// Display up to 8 curses
		if i >= maxCurses {
			// Indicate that there are more hidden curses
			str = fmt.Sprintf("%d more random curses.", len(item.Curses)-8)
			topts.GeoM.Reset()
			topts.GeoM.Translate(
				logic.ScreenWidth/2+itemCardWidth/2+48+cursePictureSize+8,
				64+float64(i*96),
			)
			text.Draw(screen, str, textFace, topts)
			break
		}
		// Name
		topts.GeoM.Reset()
		topts.GeoM.Translate(
			logic.ScreenWidth/2+itemCardWidth/2+48+cursePictureSize+8,
			64+16+float64(i*96),
		)
		text.Draw(screen, c.Name(), textFace, topts)
		// Description
		str = ""
		for i, r := range c.Description() {
			str += string(r)
			if i%55 == 0 {
				if idx := strings.LastIndex(str, " "); idx != -1 {
					str = str[:idx] + "\n" + str[idx+1:]
				}
			}
		}
		topts.LineSpacing = 1 * tinyFace.Size
		topts.GeoM.Reset()
		topts.GeoM.Translate(
			logic.ScreenWidth/2+itemCardWidth/2+48+cursePictureSize+8,
			64+48+float64(i*96),
		)
		text.Draw(screen, str, tinyFace, topts)
	}
}

func (hud *HUD) Draw(screen *ebiten.Image, ctx *HUDContext) {
	const (
		playerHPBarWidth  = 768
		playerHPBarHeight = 64
	)
	// Player HP
	vx, ix := graphics.AppendRectVerticesIndices(
		nil, nil, 0, &graphics.RectOpts{
			DstX:      logic.ScreenWidth/2 - playerHPBarWidth/2,
			DstY:      logic.ScreenHeight - 96,
			SrcX:      1,
			SrcY:      1,
			DstWidth:  playerHPBarWidth,
			DstHeight: playerHPBarHeight,
			R:         0.1,
			G:         0,
			B:         0.2,
			A:         0.7,
		},
	)
	vx, ix = graphics.AppendRectVerticesIndices(
		vx, ix, 1, &graphics.RectOpts{
			DstX:      logic.ScreenWidth/2 - playerHPBarWidth/2,
			DstY:      logic.ScreenHeight - 96,
			SrcX:      1,
			SrcY:      1,
			DstWidth:  float32(playerHPBarWidth * (ctx.PlayerHP / ctx.PlayerMaxHP)),
			DstHeight: playerHPBarHeight,
			R:         0.5,
			G:         0,
			B:         0.7,
			A:         0.7,
		},
	)

	const (
		bossHPBarWidth  = 1280
		bossHPBarHeight = 64
	)
	if ctx.StageKind == core.BossFight {
		// Boss HP
		vx, ix = graphics.AppendRectVerticesIndices(
			vx, ix, 2, &graphics.RectOpts{
				DstX:      logic.ScreenWidth/2 - bossHPBarWidth/2,
				DstY:      32,
				SrcX:      1,
				SrcY:      1,
				DstWidth:  bossHPBarWidth,
				DstHeight: bossHPBarHeight,
				R:         0.2,
				G:         0,
				B:         0.,
				A:         0.7,
			},
		)
		vx, ix = graphics.AppendRectVerticesIndices(
			vx, ix, 3, &graphics.RectOpts{
				DstX:      logic.ScreenWidth/2 - bossHPBarWidth/2,
				DstY:      32,
				SrcX:      1,
				SrcY:      1,
				DstWidth:  float32(bossHPBarWidth * (ctx.BossHP / ctx.BossMaxHP)),
				DstHeight: bossHPBarHeight,
				R:         0.7,
				G:         0,
				B:         0.,
				A:         0.7,
			},
		)
	}
	screen.DrawTriangles(vx, ix, graphics.BrushImage, nil)

	// Player HP text
	opts := &text.DrawOptions{}
	str := fmt.Sprintf("%d/%d", int(ctx.PlayerHP), int(ctx.PlayerMaxHP))
	w, h := text.Measure(str, titleFace, 0)
	opts.GeoM.Translate(
		logic.ScreenWidth/2-w/2,
		logic.ScreenHeight-96+playerHPBarHeight/2-h/2,
	)
	opts.ColorScale.Scale(0.8, 0.8, 0.8, 1)
	text.Draw(screen, str, titleFace, opts)
	// Boss HP text
	if ctx.StageKind == core.BossFight {
		str = fmt.Sprintf("%d/%d", int(ctx.BossHP), int(ctx.BossMaxHP))
		w, h = text.Measure(str, titleFace, 0)
		opts.GeoM.Reset()
		opts.GeoM.Translate(
			logic.ScreenWidth/2-w/2,
			32+bossHPBarHeight/2-h/2,
		)
		opts.ColorScale.Reset()
		opts.ColorScale.Scale(0.8, 0.8, 0.8, 1)
		text.Draw(screen, str, titleFace, opts)
	}
	// Item tooltip
	if ctx.StageKind == core.Building && ctx.TargetItem != nil {
		hud.drawItemTooltip(screen, ctx.TargetItem)
	}
	// Stage number
	str = fmt.Sprintf("Stage: %d", ctx.Stage)
	//w, h = text.Measure(str, titleFace, 0)
	opts.GeoM.Reset()
	opts.GeoM.Translate(
		logic.ScreenWidth-256,
		32,
	)
	opts.ColorScale.Scale(0.8, 0.8, 0.8, 1)
	text.Draw(screen, str, titleFace, opts)
}
