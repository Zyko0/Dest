package main

import (
	"image"
	"log"

	"github.com/Zyko0/Alapae/assets"
	"github.com/Zyko0/Alapae/core"
	"github.com/Zyko0/Alapae/input"
	"github.com/Zyko0/Alapae/logic"
	"github.com/Zyko0/Alapae/ui"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	offscreen *ebiten.Image
	game      *core.Game
	hud       *ui.HUD

	splash *ui.SplashView
	stats  *ui.Stats

	updated bool
}

func New() *Game {
	return &Game{
		offscreen: ebiten.NewImage(logic.ScreenWidth, logic.ScreenHeight),
		game: core.NewGame(core.NewCamera(
			mgl64.Vec3{0, 0, 0},
			mgl64.Vec3{0, 0, 0},
			45,
			float64(logic.ScreenWidth)/float64(logic.ScreenHeight),
		), image.Rect(0, 0, logic.ScreenWidth, logic.ScreenHeight)),
		hud: &ui.HUD{},

		splash: ui.NewSplashView(),
		stats:  ui.NewStats(),
	}
}

func (g *Game) Update() error {
	if g.splash.Active() {
		g.splash.Update()
		g.updated = true
		return nil
	}
	if g.stats.RestartGame || inpututil.IsKeyJustPressed(ebiten.KeyR) {
		//g.paused = false
		g.stats.Disable()
		g.game = core.NewGame(core.NewCamera(
			mgl64.Vec3{0, 0, 0},
			mgl64.Vec3{0, 0, 0},
			45,
			float64(logic.ScreenWidth)/float64(logic.ScreenHeight),
		), image.Rect(0, 0, logic.ScreenWidth, logic.ScreenHeight))
	}
	sctx := &ui.StatsContext{
		Title: "Pause",
		Build: g.game.Player.Core,
	}
	if g.game.Player.Dead() {
		g.stats.Enable()
		sctx.Title = "Game over"
	} else if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		if !g.stats.Active {
			g.stats.Enable()
		} else {
			g.stats.Disable()
		}
	}
	if g.stats.Active {
		input.SetLastCursor(ebiten.CursorPosition())
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
		g.stats.Update(sctx)
		g.updated = true
		return nil
	}

	if !input.EnsureCursorCaptured() {
		// TODO: don't treat input instead of returning here, but keep the
		// game running
		return nil
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
		// Stats menu
		if g.stats.Active {
			g.stats.Draw(g.offscreen)
		}
		// Splash screen
		if g.splash.Active() {
			g.splash.Draw(g.offscreen)
		}
		// Mark as drawn
		g.updated = false
	}
	screen.DrawImage(g.offscreen, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logic.ScreenWidth, logic.ScreenHeight
}

func main() {
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(core.TPS)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(logic.ScreenWidth, logic.ScreenHeight)
	ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)

	assets.SetMusic(assets.MusicMenuShop)
	assets.PlayMusic()
	if err := ebiten.RunGameWithOptions(New(), &ebiten.RunGameOptions{
		GraphicsLibrary: ebiten.GraphicsLibraryOpenGL,
	}); err != nil {
		log.Fatal("rungame:", err)
	}
}
