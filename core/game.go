package core

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"sort"

	"github.com/Zyko0/Alapae/assets"
	"github.com/Zyko0/Alapae/core/aoe"
	"github.com/Zyko0/Alapae/core/boss"
	"github.com/Zyko0/Alapae/core/building"
	"github.com/Zyko0/Alapae/core/entity"
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

	stage    int
	camera   *Camera
	player   *Player
	building *building.Phase
	boss     entity.Boss
	entities []entity.Entity
}

func NewGame(camera *Camera, resolution image.Rectangle) *Game {
	p := newPlayer()

	return &Game{
		seed:  rand.Float32(),
		floor: newFloor(),

		camera:   camera,
		player:   p,
		building: building.NewPhase(p.Core),
		entities: []entity.Entity{},
	}
}

var (
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
	g.camera.SetPosition(pos)
}

type Stage byte

const (
	Building Stage = iota
	BossFight
)

func (g *Game) Stage() Stage {
	if g.stage%2 == 0 {
		return Building
	}
	return BossFight
}

func (g *Game) InitStage() {
	if g.Stage() == Building {
		g.building.RollNew()
		g.entities = g.building.AppendEntities(g.entities)
		return
	}
	b := boss.NewSmokeMask(mgl64.Vec3{
		192 / 2,
		graphics.SpriteScale,
		192,
	})
	g.boss = b
	g.entities = append(g.entities, b)
}

func (g *Game) StageSheetImage() *ebiten.Image {
	if g.Stage() == Building {
		return assets.ItemSheetImage
	}
	return g.boss.Image()
}

func (g *Game) Update() {
	if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		g.InitStage()
	}
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
			/*&aoe.Arrow{
				X:        ArenaSize / 2,
				Y:        ArenaSize / 2,
				Size:     50,
				Rotation: 0,
			},*/
			&aoe.CircleBorder{
				X:      ArenaSize / 2,
				Y:      ArenaSize / 2,
				Radius: 45,
			},
			60, 120000,
		))
	}
	// End debug

	// Inputs and camera
	g.processInput()
	g.camera.Update()
	// Updates
	ctx := &entity.Context{
		CameraRight:     g.camera.right,
		CameraUp:        g.camera.up,
		PlayerPosition:  g.camera.position,
		PlayerDirection: g.camera.direction,
		Boss:            g.boss,
	}
	// Entities update
	var n int
	for _, p := range g.entities {
		p.Update(ctx)
		if p.Dead() {
			continue
		}
		g.entities[n] = p
		n++
	}
	g.entities = g.entities[:n]
	// Player update
	g.player.Update(ctx)
	g.entities = append(g.entities, ctx.Entities...)

	// Building phase
	if g.Stage() == Building {
		g.building.Update(ctx)
	}

	// TODO: collisions
	// Update/add floor AoE markers
	for _, m := range ctx.Markers {
		g.floor.AddMarker(m)
	}
	g.floor.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	pos := g.camera.Position()
	pvinv := g.camera.ProjectionMatrix().Mul4(g.camera.ViewMatrix()).Inv()

	// Draw entities
	ctx := &graphics.Context{
		ScreenBounds:   screen.Bounds(),
		CameraPosition: pos,
		CameraRight:    g.camera.right,
		CameraUp:       g.camera.up,
		ProjView:       g.camera.proj.Mul4(g.camera.view),
		ViewInv:        g.camera.ViewMatrix().Inv(),
	}
	var vx []ebiten.Vertex
	var ix []uint16
	// Refresh AoE markers on the floor
	var shape aoe.Shape
	if g.boss != nil {
		shape = g.boss.MarkerShape()
	} else {
		shape = g.building.MarkerShape()
	}
	g.floor.Draw(shape)
	// Draw arena scene
	vx, ix = graphics.AppendRectVerticesIndices(
		vx[:0], ix[:0], 0, &graphics.RectOpts{
			DstWidth:  float32(screen.Bounds().Dx()),
			DstHeight: float32(screen.Bounds().Dy()),
			SrcWidth:  float32(screen.Bounds().Dx()),
			SrcHeight: float32(screen.Bounds().Dy()),
		},
	)
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
	// Entities
	index := 0
	vx, ix = vx[:0], ix[:0]
	// TODO: cache the lengths
	sort.SliceStable(g.entities, func(i, j int) bool {
		li := pos.Sub(g.entities[i].Position()).LenSqr()
		lj := pos.Sub(g.entities[j].Position()).LenSqr()
		return li > lj
	})
	for _, p := range g.entities {
		vx, ix = p.AppendVerticesIndices(vx, ix, &index, ctx)
	}
	graphics.ScreenVertices(vx, screen.Bounds().Dx(), screen.Bounds().Dy())
	screen.DrawTrianglesShader(vx, ix, assets.ShaderEntity(), &ebiten.DrawTrianglesShaderOptions{
		Images: [4]*ebiten.Image{
			g.StageSheetImage(),
		},
	})

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
