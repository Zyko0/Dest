package core

import (
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
	"github.com/Zyko0/Alapae/logic"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TPS = 60

	ArenaSize = 192
)

type Game struct {
	seed   float32
	floor  *Floor
	portal *entity.Portal
	stage  int

	Player   *Player
	Building *building.Phase
	Boss     entity.Boss

	camera   *Camera
	entities []entity.Entity
}

func NewGame(camera *Camera, resolution image.Rectangle) *Game {
	p := newPlayer()

	g := &Game{
		seed:   rand.Float32(),
		floor:  newFloor(),
		portal: entity.NewPortal(),

		camera:   camera,
		Player:   p,
		Building: building.NewPhase(p.Core),
		entities: []entity.Entity{},
	}
	g.initStage()
	return g
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
		PlayerMovementSpeed*g.Player.speedMod,
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

func (g *Game) StageNumber() int {
	return g.stage
}

func (g *Game) Stage() Stage {
	if g.stage%2 == 0 {
		return Building
	}
	return BossFight
}

func (g *Game) initStage() {
	g.entities = g.entities[:0]
	if g.Stage() == Building {
		g.Building.RollNew()
		g.entities = g.Building.AppendEntities(g.entities)
		g.entities = append(g.entities, g.portal)
		g.portal.Activate()
		g.camera.SetPosition(mgl64.Vec3{
			graphics.SpriteScale,
			g.camera.position.Y(),
			logic.ArenaSize - graphics.SpriteScale,
		})
		g.camera.SetYawPitch(2.5, 0)
		// Health recover
		g.Player.Core.Health = min(
			g.Player.Core.Health+g.Player.Core.HealthPerStage,
			g.Player.Core.MaxHealth,
		)
		assets.SetMusic(assets.MusicMenuShop)
		assets.PlayMusic()
		return
	}
	b := boss.NewSmokeMask(mgl64.Vec3{
		logic.ArenaSize - graphics.SpriteScale,
		graphics.SpriteScale,
		graphics.SpriteScale,
	}, g.StageNumber())
	g.Boss = b
	g.entities = append(g.entities, b)
	g.portal.Deactivate()
	g.camera.SetPosition(mgl64.Vec3{
		graphics.SpriteScale,
		g.camera.position.Y(),
		logic.ArenaSize - graphics.SpriteScale,
	})
	g.camera.SetYawPitch(2.5, 0)
	assets.SetMusic(assets.MusicBoss0)
	assets.PlayMusic()
}

func (g *Game) nextStage() {
	g.stage++
	g.seed = rand.Float32()
	g.initStage()
}

func (g *Game) StageSheetImage() *ebiten.Image {
	if g.Stage() == Building {
		return assets.ItemSheetImage
	}
	return g.Boss.Image()
}

func (g *Game) Update() {
	// Inputs and camera
	g.processInput()
	g.camera.Update()
	// Entity collisions
	g.handleCollisions()
	// Updates
	ctx := &entity.Context{
		CameraRight:     g.camera.right,
		CameraUp:        g.camera.up,
		PlayerPosition:  g.camera.position,
		PlayerDirection: g.camera.direction,
		Boss:            g.Boss,
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
	g.Player.Update(ctx)
	g.entities = append(g.entities, ctx.Entities...)

	// Building phase
	g.Player.LeftHand.Glow, g.Player.RightHand.Glow = 0, 0
	switch g.Stage() {
	case Building:
		g.Building.Update(ctx)
		// Hand glowing
		if g.Building.Target != nil {
			switch g.Building.Target.Item.HandSide {
			case building.BothHand:
				g.Player.LeftHand.Glow, g.Player.RightHand.Glow = 1, 1
			case building.RightHand:
				g.Player.RightHand.Glow = 1
			case building.LeftHand:
				g.Player.LeftHand.Glow = 1
			}
		}
	case BossFight:
		if g.portal.Dead() && (g.Boss == nil || g.Boss.Dead()) {
			g.entities = append(g.entities, g.portal)
			g.portal.Activate()
		}
	}
	// Interact
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		switch {
		case g.Stage() == Building && g.Building.Target != nil:
			item := g.Building.Target.Item
			g.Building.Pick()
			item.RegisterMod(g.Player.Core, g.Building)
			g.Player.Core.Update()
			assets.PlayBonusPickup()
		case !g.portal.Dead() && g.portal.Targeted():
			g.nextStage()
			assets.PlayPortal()
			return
		}
	}

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
	if g.Boss != nil && !g.Boss.Dead() {
		shape = g.Boss.MarkerShape()
	} else {
		shape = g.Building.MarkerShape()
	}
	if shape == nil && !g.portal.Dead() {
		shape = g.portal.MarkerShape()
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
	// TODO: cache the lengths?
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
	g.Player.DrawHands(screen, ctx)
	// Crosshair
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(
		float64(screen.Bounds().Dx())/2-float64(assets.CursorImage.Bounds().Dx())/2,
		float64(screen.Bounds().Dy())/2-float64(assets.CursorImage.Bounds().Dy())/2,
	)
	screen.DrawImage(assets.CursorImage, opts)
}
