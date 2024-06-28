package core

import (
	"github.com/Zyko0/Alapae/assets"
	"github.com/Zyko0/Alapae/core/building"
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
	InvulnTicks     = 20
	InvulnCD        = InvulnTicks
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
	shooting
)

type Player struct {
	status   status
	speedMod float64
	active   *state
	cooldown *state

	Core       *building.Core
	ActiveHand *hand.Hand
	RightHand  *hand.Hand
	LeftHand   *hand.Hand
}

func newPlayer() *Player {
	p := &Player{
		status:    idle,
		Core:      building.NewCore(),
		RightHand: hand.New(hand.Right),
		LeftHand:  hand.New(hand.Left),
		active:    &state{},
		cooldown:  &state{},
	}
	p.ActiveHand = p.RightHand
	p.resetModifiers()

	return p
}

func (p *Player) resetModifiers() {
	p.speedMod = 1
	// TODO: ?
}

func (p *Player) Update(ctx *entity.Context) {
	var hands []*hand.Hand
	if p.Core.Synced {
		hands = []*hand.Hand{p.RightHand, p.LeftHand}
	} else {
		hands = []*hand.Hand{p.ActiveHand}
	}

	// Shooting
	if p.status == idle && p.cooldown.Shooting == 0 && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p.active.Shooting = ShootingTicks
		p.cooldown.Shooting = ShootingCD
		p.status = shooting
		// Swap hand if necessary
		for _, h := range hands {
			h.Anim = h.ShootAnimation(p.Core.Hand(h.Side).Weapon).NewInstance(h, false)
			if h == p.RightHand {
				p.ActiveHand = p.LeftHand
			} else {
				p.ActiveHand = p.RightHand
			}
		}
	}
	// Dashing
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if p.active.Dashing == 0 && p.cooldown.Dashing == 0 {
			p.active.Dashing = DashingTicks
			p.cooldown.Dashing = DashingCD
			p.active.Invuln = InvulnTicks
			p.cooldown.Invuln = InvulnCD
		}
	}
	// Hand animation test // TODO:
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {

	}

	// Update current effects
	p.active.Update()
	p.cooldown.Update()
	p.resetModifiers()
	// Ending status
	for _, h := range hands {
		switch {
		case p.status == shooting && p.active.Shooting == 0:
			off := ctx.PlayerDirection.Mul(0.5 * graphics.SpriteScale)
			off = off.Add(ctx.CameraRight.Mul(0.75 * h.ShotRightCoeff()))
			off = off.Sub(ctx.CameraUp.Mul(0.5))
			// Shoot a projectile
			data := p.Core.Projectile(h.Side)
			ctx.Entities = append(ctx.Entities, entity.NewProjectile(
				ctx.PlayerPosition.Add(off),
				ctx.PlayerDirection,
				data.Radius,
				data.Speed,
				data.Damage,
				entity.TeamAlly,
				data.ColorIn,
				data.ColorOut,
				data.Alpha,
				data.MaxDuration,
				data.Resistance,
			))
			// If not hand-synced, terminate the other hand's animation
			//if !p.Core.Synced {
			if h == p.RightHand {
				p.LeftHand.Anim = nil
			} else {
				p.RightHand.Anim = nil
			}
			//}
		}
	}
	// New states
	if p.active.Dashing > 0 {
		p.speedMod = DashingSpeedMod
	}
	switch {
	case p.active.Shooting > 0:
		p.status = shooting
	default:
		p.status = idle
		if p.RightHand.Anim == nil {
			p.RightHand.Anim = hand.AnimationIdle.NewInstance(p.RightHand, false)
		}
		if p.LeftHand.Anim == nil {
			p.LeftHand.Anim = hand.AnimationIdle.NewInstance(p.LeftHand, true)
		}
	}
	// Update animations
	p.RightHand.Anim.Update(p.RightHand)
	p.LeftHand.Anim.Update(p.LeftHand)
}

func (p *Player) TakeHit(dmg float64) {
	if dmg > 0 && p.active.Invuln <= 0 {
		// TODO: play sfx
		p.Core.Health = max(p.Core.Health-dmg, 0)
		p.active.Invuln = InvulnTicks
		p.cooldown.Invuln = InvulnCD
	}
}

func (p *Player) Dead() bool {
	return p.Core.Health <= 0
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
