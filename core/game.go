package core

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/Zyko0/Alapae/assets"
	"github.com/Zyko0/Alapae/core/aoe"
	"github.com/Zyko0/Alapae/graphics"
	"github.com/Zyko0/Alapae/input"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TPS = 60

	ArenaSize = 192
)

type Game struct {
	seed  float32
	floor *Floor

	camera   *Camera
	player   *Player
	boss     Boss
	entities []Entity
}

func NewGame(camera *Camera) *Game {
	return &Game{
		seed:  rand.Float32(),
		floor: newFloor(),

		camera: camera,
		player: newPlayer(),
	}
}

var (
	Grounded   bool
	PlayerSize = mgl64.Vec3{1, 2, 1}
)

func (g *Game) processInput() {
	// Mouse
	yaw, pitch := g.camera.YawPitch()
	yawoff, pitchoff := input.ProcessMouseMovement()
	yaw += yawoff
	pitch += pitchoff
	if pitch > math.Pi/2 {
		pitch = math.Pi / 2
	}
	if pitch < -math.Pi/2 {
		pitch = -math.Pi / 2
	}
	g.camera.SetYawPitch(yaw, pitch)
	// Keyboard
	pos := input.ProcessKeyboard(
		g.camera.Position(),
		g.camera.Direction(),
		g.camera.Right(),
		PlayerMovementSpeed*g.player.SpeedMod,
	)
	pos[0] = max(min(pos[0], ArenaSize-PlayerSize.X()/2), 0)
	pos[2] = max(min(pos[2], ArenaSize-PlayerSize.Z()/2), 0)
	if pos[1] < 1+PlayerSize.Y()/2 {
		pos[1] = 1 + PlayerSize.Y()/2
		Grounded = true
	}
	g.camera.SetPosition(pos)
}

func (g *Game) Update() {
	// TODO: Debug
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.floor.AddMarker(aoe.NewMarker(
			/*&aoe.Circle{
				X:      ArenaSize / 2,
				Y:      ArenaSize / 2,
				Radius: 45,
			},*/
			/*&aoe.XCross{
				Size:         ArenaSize,
				Radius:       0.25,
				Rotation:     0,
				RotationIncr: 0.025,
			},*/
			&aoe.Arrow{
				X:        ArenaSize / 2,
				Y:        ArenaSize / 2,
				Size:     50,
				Rotation: 0,
			},
			60, 120000,
		))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.seed = rand.Float32()
	}
	// End debug

	// Inputs and camera
	g.processInput()
	g.camera.Update()
	// Players actions
	ctx := &PlayerContext{
		PlayerPosition:  g.camera.position,
		PlayerDirection: g.camera.direction,
	}
	g.player.Update(ctx)
	// TODO: g.boss.Update(ctx)
	g.entities = append(g.entities, ctx.Projectiles...)
	// Entities update
	var n int
	for _, e := range g.entities {
		e.Update()
		if e.Dead() {
			continue
		}
		g.entities[n] = e
		n++
	}
	g.entities = g.entities[:n]
	// TODO: collisions
	// Update floor AoE markers
	g.floor.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Refresh AoE markers on the floor
	g.floor.Draw()
	// Draw arena scene
	vx, ix := graphics.AppendRectVerticesIndices(
		nil, nil, 0, &graphics.RectOpts{
			DstWidth:  float32(screen.Bounds().Dx()),
			DstHeight: float32(screen.Bounds().Dy()),
			SrcWidth:  float32(screen.Bounds().Dx()),
			SrcHeight: float32(screen.Bounds().Dy()),
		},
	)
	pos := g.camera.Position()
	pvinv := g.camera.ProjectionMatrix().Mul4(g.camera.ViewMatrix()).Inv()
	screen.DrawTrianglesShader(vx, ix, assets.ShaderArena(), &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]any{
			"Seed":                   g.seed,
			"CameraPosition":         pos[:],
			"CameraPVMatrixInv":      pvinv[:],
			"MarkerResolutionFactor": float32(markerResolutionFactor),
		},
		Images: [4]*ebiten.Image{
			g.floor.Image,
		},
	})
	// Draw entities
	ctx := &graphics.Context{
		CameraPosition: pos,
		CameraRight:    g.camera.right,
		CameraUp:       g.camera.up,
		ProjView:       g.camera.proj.Mul4(g.camera.view),
		ViewInv:        g.camera.ViewMatrix().Inv(),
	}
	vx, ix = vx[:0], ix[:0]
	index := 0
	// TODO: sort them by Z
	for _, e := range g.entities {
		vx, ix = e.AppendVerticesIndices(vx, ix, &index, ctx)
	}
	graphics.ScreenVertices(vx, screen.Bounds().Dx(), screen.Bounds().Dy())
	//fmt.Println("len entities:", len(g.entities), "len vx:", len(vx), "len ix:", len(ix))
	// TODO: add src images corresponding to billboard sprites
	screen.DrawTrianglesShader(vx, ix, assets.ShaderEntity(), nil)

	// Player hands
	g.player.DrawHands(screen, ctx)
	// TODO: this is debug
	screen.DrawImage(g.floor.Image, nil)
	// Crosshair
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(
		float64(screen.Bounds().Dx())/2-float64(assets.CursorImage.Bounds().Dx())/2,
		float64(screen.Bounds().Dy())/2-float64(assets.CursorImage.Bounds().Dy())/2,
	)
	screen.DrawImage(assets.CursorImage, opts)

	// Debug
	y, p := g.camera.YawPitch()
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("TPS: %0.2f - FPS %.02f - y %.02f p %.02f - pos %v - dir %v",
			ebiten.ActualTPS(),
			ebiten.ActualFPS(),
			y, p, g.camera.Position(), g.camera.Direction(),
		),
	)
}
