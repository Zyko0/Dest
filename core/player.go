package core

import (
	"image/color"

	"github.com/Zyko0/Alapae/assets"
	"github.com/Zyko0/Alapae/core/entity"
	"github.com/Zyko0/Alapae/core/hand"
	"github.com/Zyko0/Alapae/graphics"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	PlayerMovementSpeed = 1.

	ShootingTicks   = 15
	ShootingCD      = 30
	DashingTicks    = 10
	DashingCD       = 60
	DashingSpeedMod = 3
)

type state struct {
	Invuln   int
	Shooting int
	Dashing  int
}

func (s *state) Update() {
	s.Invuln = max(s.Invuln-1, 0)
	s.Shooting = max(s.Shooting-1, 0)
	s.Dashing = max(s.Dashing-1, 0)
}

type status byte

const (
	idle status = iota
	moving
	dashing
	shooting
)

type Player struct {
	Health    int
	MaxHealth int

	Status status

	ActiveHand *hand.Hand
	RightHand  *hand.Hand
	LeftHand   *hand.Hand
	Active     *state
	Cooldown   *state

	SpeedMod           float64
	ProjectileSpeedMod float64
	// TODO: hands
	// TODO: powerups
	// TODO: curses
}

func newPlayer() *Player {
	p := &Player{
		Health:    100,
		MaxHealth: 100,

		RightHand: hand.New(hand.Right),
		LeftHand:  hand.New(hand.Left),
		Active:    &state{},
		Cooldown:  &state{},
	}
	p.ActiveHand = p.RightHand
	p.resetModifiers()

	return p
}

func (p *Player) resetModifiers() {
	p.SpeedMod = 1
	p.ProjectileSpeedMod = 2
}

func (p *Player) Update(ctx *entity.Context) {
	// Shooting
	if p.Status == idle && p.Cooldown.Shooting == 0 && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p.Active.Shooting = ShootingTicks
		p.Cooldown.Shooting = ShootingCD
		p.Status = shooting
		// TODO: swap hands
		if p.ActiveHand == p.RightHand {
			p.RightHand.Anim = hand.AnimationShootFinger.NewInstance(p.RightHand, false)
			p.ActiveHand = p.LeftHand
		} else {
			p.LeftHand.Anim = hand.AnimationShootFinger.NewInstance(p.LeftHand, false)
			p.ActiveHand = p.RightHand
		}
	}
	// Dashing
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if p.Active.Dashing == 0 && p.Cooldown.Dashing == 0 {
			p.Active.Dashing = DashingTicks
			p.Cooldown.Dashing = DashingCD
		}
	}
	// Hand animation test
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		/*ctx.Entities = append(ctx.Entities, entity.NewComet(
			mgl64.Vec3{192 / 2, 100, 192 / 2},
			1, 2,
		))*/
	}

	// Update current effects
	p.Active.Update()
	p.Cooldown.Update()
	p.resetModifiers()
	// Ending status
	switch {
	case p.Status == shooting && p.Active.Shooting == 0:
		off := ctx.PlayerDirection.Mul(0.5 * graphics.SpriteScale)
		off = off.Add(ctx.CameraRight.Mul(0.75 * p.ActiveHand.ShotRightCoeff()))
		off = off.Sub(ctx.CameraUp.Mul(0.5))

		ctx.Entities = append(ctx.Entities, entity.NewProjectile(
			ctx.PlayerPosition.Add(off),
			ctx.PlayerDirection,
			0.1,
			p.ProjectileSpeedMod,
			color.RGBA{0, 0, 255, 255},
			color.White,
		))
		if p.ActiveHand == p.RightHand {
			p.LeftHand.Anim = nil
		} else {
			p.RightHand.Anim = nil
		}
	}
	// New states
	switch {
	case p.Active.Dashing > 0:
		p.SpeedMod = DashingSpeedMod
		p.Status = dashing
	case p.Active.Shooting > 0:
		p.Status = shooting
	default:
		p.Status = idle
		if p.RightHand.Anim == nil {
			p.RightHand.Anim = hand.AnimationIdle.NewInstance(p.RightHand, false)
		}
		if p.LeftHand.Anim == nil {
			p.LeftHand.Anim = hand.AnimationIdle.NewInstance(p.LeftHand, true)
		}
	}
	// Update animations
	if p.RightHand.Anim != nil {
		p.RightHand.Anim.Update(p.RightHand)
	}
	if p.LeftHand.Anim != nil {
		p.LeftHand.Anim.Update(p.LeftHand)
	}
}

func (p *Player) DrawHands(screen *ebiten.Image, ctx *graphics.Context) {
	const (
		size = 512
		yoff = size / 8
	)
	// Right
	vx, ix := graphics.AppendRectVerticesIndices(nil, nil, 0, &graphics.RectOpts{
		DstX:      float32(screen.Bounds().Dx()) / 2,
		DstY:      float32(screen.Bounds().Dy()) - size - yoff,
		SrcX:      -1,
		SrcY:      -1,
		DstWidth:  size,
		DstHeight: size,
		SrcWidth:  2,
		SrcHeight: 2,
		R:         p.RightHand.Glow,
		G:         0,
		B:         0,
		A:         0,
	})
	var data []float32
	data = p.RightHand.AppendData(data)
	screen.DrawTrianglesShader(vx, ix, assets.ShaderHands(), &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]any{ // TODO: might be useless uniforms
			"Rotation": []float32{
				float32(p.RightHand.Rotation[0]),
				float32(p.RightHand.Rotation[1]),
				float32(p.RightHand.Rotation[2]),
			},
			"Fingers": data,
		},
	})
	// Left
	vx, ix = graphics.AppendRectVerticesIndices(vx[:0], ix[:0], 0, &graphics.RectOpts{
		DstX:      float32(screen.Bounds().Dx())/2 - size,
		DstY:      float32(screen.Bounds().Dy()) - size - yoff,
		SrcX:      1,
		SrcY:      -1,
		DstWidth:  size,
		DstHeight: size,
		SrcWidth:  -2,
		SrcHeight: 2,
		R:         p.LeftHand.Glow,
		G:         0,
		B:         0,
		A:         0,
	})
	data = p.LeftHand.AppendData(data[:0])
	screen.DrawTrianglesShader(vx, ix, assets.ShaderHands(), &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]any{ // TODO: might be useless uniforms
			"Rotation": []float32{
				float32(p.LeftHand.Rotation[0]),
				float32(p.LeftHand.Rotation[1]),
				float32(p.LeftHand.Rotation[2]),
			},
			"Fingers": data,
		},
	})
}
