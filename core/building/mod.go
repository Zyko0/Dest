package building

import (
	"math/rand"
	"slices"

	"github.com/Zyko0/Alapae/core/hand"
)

type Mod struct {
	def *Definition

	Stacks int
}

func (m *Mod) Init(c *Core, hm *HandModifiers, p *Phase) {
	switch m.def.ID {
	case Prayer:
		if len(hm.Curses) > 0 {
			i := rand.Intn(len(hm.Curses))
			hm.Curses = slices.Delete(hm.Curses, i, i+1)
		}
		c.Health = max(min(c.Health+c.MaxHealth*0.2, c.MaxHealth), 1)
	case Highroll:
		p.RollExisting(p.Objects, false)
	case Striker, Gambler:
		m.Stacks = c.AttackSpeedStacks
		c.AttackSpeedStacks = 0
	case Dual_Prayer:
		if len(c.right.Curses) > 0 {
			r := rand.Intn(len(c.right.Curses))
			c.right.Curses = slices.Delete(c.right.Curses, r, r+1)
		}
		if len(c.left.Curses) > 0 {
			l := rand.Intn(len(c.left.Curses))
			c.left.Curses = slices.Delete(c.left.Curses, l, l+1)
		}
		c.Health = max(min(c.Health+c.MaxHealth*0.2, c.MaxHealth), 1)
	case Change_of_mind:
		c.right.Bonuses, c.left.Bonuses = c.left.Bonuses, c.right.Bonuses
	case Mimic:
		var other *HandModifiers
		if hm == c.left {
			other = c.right
		} else {
			other = c.left
		}
		hm.Bonuses, hm.Curses = hm.Bonuses[:0], hm.Curses[:0]
		for _, m := range other.Bonuses {
			hm.Bonuses = append(hm.Bonuses, &Mod{
				def:    m.def,
				Stacks: m.Stacks,
			})
		}
		for _, m := range other.Curses {
			hm.Curses = append(hm.Curses, &Mod{
				def:    m.def,
				Stacks: m.Stacks,
			})
		}
	case Trap:
		c.Health = max(c.Health-max(c.MaxHealth*0.1, 1), 0)
	case Lowroll:
		p.RollExisting(p.Objects, false)
	case Sabotage:
		p.RollExtraCurses(10)
	case Procrastination:
		p.RegisterExtraCurse()
	}
	// One apply
	m.Apply(c, hm)
}

func (m *Mod) Apply(c *Core, hm *HandModifiers) {
	switch m.def.ID {
	case Damage_up:
		hm.Damage += 5
	case Critical_chance:
		if hm.CritChance == 1 {
			hm.CritDamage += 0.05
		} else {
			hm.CritChance += 0.05
		}
	case Prayer:
	case Luck:
		c.Luck = min(c.Luck+0.025, 1)
	case Highroll:
	case Dual_damage_up:
		c.right.Damage += 2.5 * float64(m.Stacks)
		c.left.Damage += 2.5 * float64(m.Stacks)
	case Attack_speed_up:
		c.AttackSpeedStacks = min(c.AttackSpeedStacks+1, 6)
	case Striker:
		hm.Damage += float64(5 * m.Stacks)
	case Gambler:
		for i := 0; i < m.Stacks; i++ {
			if hm.CritChance == 1 {
				hm.CritDamage += 0.05
			} else {
				hm.CritChance += 0.05
			}
		}
	case Dual_Prayer:
	case Dual_damage_way_up:
		c.right.Damage += 5 * float64(m.Stacks)
		c.left.Damage += 5 * float64(m.Stacks)
	case Critical_damage:
		c.right.CritDamage += 0.1
		c.left.CritDamage += 0.1
	case Curse_advantage_ex:
		var count int
		if hm == c.left {
			count = len(c.right.Curses)
		} else {
			count = len(c.left.Curses)
		}
		hm.Damage += float64(5 * count)
	case Curse_advantage:
		hm.Damage += float64(5 * len(hm.Curses))
	case Survivor:
		hm.Damage += float64(5 * ((100 - c.Health) / 5))
	case Dual_critical_chance:
		if c.right.CritChance == 1 {
			c.right.CritDamage += 0.05
		} else {
			c.right.CritChance += 0.05
		}
		if c.left.CritChance == 1 {
			c.left.CritDamage += 0.05
		} else {
			c.left.CritChance += 0.05
		}
	case Pistol:
		hm.Weapon = hand.WeaponPistol
	case Extra_shot:
		hm.ProjectileCount++
	case Dual_curse_advantage:
		hm.Damage += float64(5 * (len(c.right.Curses) + len(c.left.Curses) + len(c.Curses)))
	case Change_of_mind:
	case Homing:
		hm.Homing = true
	case Sync:
		c.Synced = true
	case Mimic:
	case Relaxed:
		c.AttackSpeedStacks -= 1
	case Clumsy:
		hm.CritDamage = max(hm.CritDamage-0.2, 0)
	case Scared:
		c.MaxHealth = max(5, c.MaxHealth-5)
		c.Health = min(c.Health, c.MaxHealth)
	case Inaccurate:
		hm.Accuracy -= 0.125
	case Heavy:
		hm.ProjectileSpeed -= 0.25
	case Trap:
	case Lowroll:
	case Sabotage:
	case Delicate:
		hm.Damage -= 10
	case Love:
		hm.InverseKnockback += 0.05
	case Procrastination:
	case Rest:
		c.HealthPerStage++
	}
}

type HandModifiers struct {
	Weapon hand.Weapon

	Damage                 float64
	CritChance             float64
	CritDamage             float64
	Accuracy               float64
	ProjectileSpeed        float64
	ProjectileCount        int
	InverseKnockback       float64
	CurseDamageCurrentHand float64
	CurseDamageOtherHand   float64
	Homing                 bool

	Bonuses []*Mod
	Curses  []*Mod
}

func (hm *HandModifiers) reset() {
	hm.Weapon = hand.WeaponFinger
	hm.Damage = 5
	hm.CritChance = 0.05
	hm.CritDamage = 2
	hm.Accuracy = 1
	hm.ProjectileSpeed = 2
	hm.ProjectileCount = 1
	hm.InverseKnockback = 0
	hm.CurseDamageCurrentHand = 0
	hm.CurseDamageOtherHand = 0
	hm.Homing = false
}

func (hm *HandModifiers) allowedStacks(id int) int {
	if impls[id].MaxStacks == InfiniteStacks {
		return 99999
	}
	var active *Mod
	for _, m := range hm.Bonuses {
		if m.def.ID == id {
			active = m
			goto stacks
		}
	}
	for _, m := range hm.Curses {
		if m.def.ID == id {
			active = m
			goto stacks
		}
	}

stacks:
	if active == nil {
		return impls[id].MaxStacks
	}

	return max(impls[id].MaxStacks-active.Stacks, 0)
}

func newHandModifiers() *HandModifiers {
	hm := &HandModifiers{}
	hm.reset()
	return hm
}

type Core struct {
	Health    float64
	MaxHealth float64

	AttackSpeedStacks int
	SurvivorStacks    int
	Synced            bool
	Luck              float64
	HealthPerStage    float64

	Bonuses []*Mod
	Curses  []*Mod

	right *HandModifiers
	left  *HandModifiers
}

func (c *Core) allowedStacks(id int) int {
	if impls[id].MaxStacks == InfiniteStacks {
		return 99999
	}
	var active *Mod
	for _, m := range c.Bonuses {
		if m.def.ID == id {
			active = m
			goto stacks
		}
	}
	for _, m := range c.Curses {
		if m.def.ID == id {
			active = m
			goto stacks
		}
	}

stacks:
	if active == nil {
		return impls[id].MaxStacks
	}

	return max(impls[id].MaxStacks-active.Stacks, 0)
}

func (c *Core) reset() {
	c.MaxHealth = 100
	c.AttackSpeedStacks = 0
	c.SurvivorStacks = 0
	c.Synced = false
	c.Luck = 0
	c.HealthPerStage = 0
	c.right.reset()
	c.left.reset()
}

func NewCore() *Core {
	c := &Core{
		Health: 100,
		right:  newHandModifiers(),
		left:   newHandModifiers(),
	}
	c.reset()
	return c
}

func (c *Core) Hand(side hand.Side) *HandModifiers {
	if side == hand.Right {
		return c.right
	}
	return c.left
}

func (c *Core) Update() {
	c.reset()
	for _, i := range c.Bonuses {
		i.Apply(c, nil)
	}
	for _, i := range c.Curses {
		i.Apply(c, nil)
	}
	for _, i := range c.right.Bonuses {
		i.Apply(c, c.right)
	}
	for _, i := range c.right.Curses {
		i.Apply(c, c.right)
	}
	for _, i := range c.left.Bonuses {
		i.Apply(c, c.left)
	}
	for _, i := range c.left.Curses {
		i.Apply(c, c.left)
	}
}
