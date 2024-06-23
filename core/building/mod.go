package building

import "github.com/Zyko0/Alapae/core/hand"

type Mod struct {
	Kind Kind
	Over bool
}

type HandStatistics struct {
	Weapon hand.Weapon

	Damage                 float64
	CritChance             float64
	CritDamage             float64
	Accuracy               float64
	ProjectileSpeed        float64
	InverseKnockback       float64
	CurseDamageCurrentHand float64
	CurseDamageOtherHand   float64
	Homing                 bool

	NextItemRandom bool

	Bonuses []*Mod
	Curses  []*Mod
}

func newHandStatistics() *HandStatistics {
	return &HandStatistics{
		Weapon:                 hand.WeaponFinger,
		Damage:                 5,
		CritChance:             0.05,
		CritDamage:             2,
		Accuracy:               1,
		ProjectileSpeed:        2,
		InverseKnockback:       0,
		CurseDamageCurrentHand: 0,
		CurseDamageOtherHand:   0,
		Homing:                 false,
		NextItemRandom:         false,
	}
}

type Statistics struct {
	Health    float64
	MaxHealth float64

	AttackSpeedStacks int
	SurvivorStacks    int
	Synced            bool
	Luck              float64

	right   *HandStatistics
	left    *HandStatistics
	bonuses []*Mod
	curses  []*Mod
}

func NewStatistics() *Statistics {
	return &Statistics{
		Health:            100,
		MaxHealth:         100,
		AttackSpeedStacks: 0,
		SurvivorStacks:    0,
		Synced:            false,
		Luck:              0,
		right:             newHandStatistics(),
		left:              newHandStatistics(),
	}
}

func (s *Statistics) Hand(side hand.Side) *HandStatistics {
	if side == hand.Right {
		return s.right
	}
	return s.left
}

func (s *Statistics) RegisterMod() {
	
}
