package main

import (
	"image"
	"log"

	"github.com/Zyko0/Alapae/core"
	"github.com/Zyko0/Alapae/input"
	"github.com/Zyko0/Alapae/ui"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
)

var (
	PlayerSize = mgl64.Vec3{1, 2, 1}
	Grounded   bool
)

type Game struct {
	offscreen *ebiten.Image
	game      *core.Game
	hud       *ui.HUD

	updated bool
}

func New() *Game {
	return &Game{
		offscreen: ebiten.NewImage(ScreenWidth, ScreenHeight),
		game: core.NewGame(core.NewCamera(
			mgl64.Vec3{0, 0, 0},
			mgl64.Vec3{0, 0, 0},
			45,
			float64(ScreenWidth)/float64(ScreenHeight),
		), image.Rect(0, 0, ScreenWidth, ScreenHeight)),
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if !input.EnsureCursorCaptured() {
		// TODO: don't treat input instead of returning here, but keep the
		// game running
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		// TODO: make this less brutal
		g.game = core.NewGame(core.NewCamera(
			mgl64.Vec3{0, 0, 0},
			mgl64.Vec3{0, 0, 0},
			45,
			float64(ScreenWidth)/float64(ScreenHeight),
		), image.Rect(0, 0, ScreenWidth, ScreenHeight))
	}

	g.game.Update()
	g.updated = true

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.updated {
		g.game.Draw(g.offscreen)
		// HUD
		hudCtx := &ui.HUDContext{
			Stage:       g.game.StageNumber(),
			StageKind:   g.game.Stage(),
			PlayerHP:    g.game.Player.Core.Health,
			PlayerMaxHP: g.game.Player.Core.MaxHealth,
		}
		switch hudCtx.StageKind {
		case core.Building:
			if g.game.Building.Target != nil {
				hudCtx.TargetItem = g.game.Building.Target.Item
			}
		case core.BossFight:
			hudCtx.BossHP = g.game.Boss.Health()
			hudCtx.BossMaxHP = g.game.Boss.MaxHealth()
		}
		g.hud.Draw(g.offscreen, hudCtx)
		g.updated = false
	}
	screen.DrawImage(g.offscreen, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	// Note: Force opengl
	ebiten.SetVsyncEnabled(false)
	ebiten.SetTPS(core.TPS)
	ebiten.SetFullscreen(true) // TODO: true
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

	if err := ebiten.RunGameWithOptions(New(), &ebiten.RunGameOptions{
		GraphicsLibrary: ebiten.GraphicsLibraryOpenGL,
	}); err != nil {
		log.Fatal("rungame:", err)
	}
}
