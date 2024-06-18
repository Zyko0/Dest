package core

import (
	"github.com/Zyko0/Alapae/assets"
	"github.com/Zyko0/Alapae/core/entity"
	"github.com/Zyko0/Alapae/core/hand"
	"github.com/Zyko0/Alapae/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	PlayerMovementSpeed = 1.

	ShootingTicks   = 20
	ShootingCD      = 30
	DashingTicks    = 15
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
	ActiveAnim *hand.AnimationInstance
	RightHand  *hand.Hand
	RightAnim  *hand.AnimationInstance
	LeftHand   *hand.Hand
	LeftAnim   *hand.AnimationInstance
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

		RightHand: hand.New(),
		LeftHand:  hand.New(),
		Active:    &state{},
		Cooldown:  &state{},
	}
	p.ActiveHand = p.RightHand
	p.resetModifiers()

	return p
}

type PlayerContext struct {
	PlayerPosition  mgl64.Vec3
	PlayerDirection mgl64.Vec3

	Projectiles []Entity
}

func (p *Player) resetModifiers() {
	p.SpeedMod = 1
	p.ProjectileSpeedMod = 4
}

func (p *Player) Update(ctx *PlayerContext) {
	// Shooting
	if p.Status == idle && p.Cooldown.Shooting == 0 && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p.Active.Shooting = ShootingTicks
		p.Cooldown.Shooting = ShootingCD
		p.Status = shooting
		// TODO: swap hands
		if p.ActiveHand == p.RightHand {
			p.RightAnim = hand.AnimationShootPistol.NewInstance(p.RightHand, false)
			p.ActiveHand = p.LeftHand
			p.ActiveAnim = p.LeftAnim
		} else {
			p.LeftAnim = hand.AnimationShootPistol.NewInstance(p.LeftHand, false)
			p.ActiveHand = p.RightHand
			p.ActiveAnim = p.RightAnim
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
		//p.RightAnim = hand.TestAnimation.Value().NewInstance(p.RightHand, false)
		//p.LeftAnim = hand.TestAnimation.Value().NewInstance(p.LeftHand, true)
	}

	// Update current effects
	p.Active.Update()
	p.Cooldown.Update()
	p.resetModifiers()
	// Ending status
	switch {
	case p.Status == shooting && p.Active.Shooting == 0:
		ctx.Projectiles = append(ctx.Projectiles, entity.NewProjectile(
			// TODO: handle right/left hand for the projectile to originate from
			ctx.PlayerPosition.Add(ctx.PlayerDirection.Mul(20)),
			ctx.PlayerDirection,
			5,
			p.ProjectileSpeedMod,
		))
		if p.ActiveHand == p.RightHand {
			p.LeftAnim = nil
		} else {
			p.RightAnim = nil
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
		if p.RightAnim == nil {
			p.RightAnim = hand.AnimationIdle.NewInstance(p.RightHand, false)
		}
		if p.LeftAnim == nil {
			p.LeftAnim = hand.AnimationIdle.NewInstance(p.LeftHand, true)
		}
	}
	// Update animations
	if p.RightAnim != nil {
		p.RightAnim.Update(p.RightHand)
	}
	if p.LeftAnim != nil {
		p.LeftAnim.Update(p.LeftHand)
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
		R:         0,
		G:         0,
		B:         0,
		A:         0,
	})
	var data []float32
	// Right
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
		R:         0,
		G:         0,
		B:         0,
		A:         0,
	})
	// Right
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
