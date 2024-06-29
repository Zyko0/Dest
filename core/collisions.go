package core

import (
	"github.com/Zyko0/Alapae/core/entity"
	"github.com/Zyko0/Alapae/graphics"
)

func (g *Game) handleCollisions() {
	if g.Stage() == Building {
		return
	}
	if g.Boss == nil || g.Boss.Dead() {
		return
	}

	const (
		playerRadiusSq = graphics.SpriteScale * graphics.SpriteScale
	)
	bossRadiusSq := g.Boss.Radius() * graphics.SpriteScale
	bossRadiusSq *= bossRadiusSq
	bossPullAmount := 0.
	var playerTest bool
	for _, e := range g.entities {
		// Player
		if !playerTest && e.Team() == entity.TeamEnemy {
			lenSq := g.camera.position.Sub(e.Position()).LenSqr()
			r := e.Radius() * graphics.SpriteScale
			lenSq -= (playerRadiusSq + r*r)
			if lenSq < 0 {
				g.Player.TakeHit(e.Damage())
				e.TakeHit(0)
				playerTest = true
			}
		}
		// Boss
		if e.Team() == entity.TeamAlly {
			lenSq := g.Boss.Position().Sub(e.Position()).LenSqr()
			r := e.Radius() * graphics.SpriteScale
			lenSq -= (bossRadiusSq + r*r)
			if lenSq < 0 {
				g.Boss.TakeHit(e.Damage())
				if proj, ok := e.(*entity.Projectile); ok {
					bossPullAmount += proj.Pull()
				}
				e.TakeHit(0)
			}
		}
	}
	// Apply inverse knockback on Boss if any
	if bossPullAmount > 0 {
		pos := g.Boss.Position()
		pos = pos.Sub(g.camera.direction.Mul(bossPullAmount))
		g.Boss.SetPosition(pos)
	}
}
