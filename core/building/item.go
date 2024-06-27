package building

import (
	"fmt"
	"image"
)

const InfiniteStacks = -1

type Kind byte

const (
	Common Kind = iota
	Uncommon
	Rare
	Epic
	Legendary
	Cursed
	Dynamic
)

type HandConstraint byte

const (
	SingleHand HandConstraint = iota
	BothHands
	None
)

type Definition struct {
	ID          int
	Kind        Kind
	Hand        HandConstraint
	Name        string
	Description string
	MaxStacks   int
	Rect        image.Rectangle
}

type HandSide byte

const (
	RightHand HandSide = iota
	LeftHand
	BothHand
	NoHand
)

func (hs HandSide) String() string {
	switch hs {
	case RightHand:
		return "right hand"
	case LeftHand:
		return "left hand"
	case BothHand:
		return "both hands"
	}
	return ""
}

type Item struct {
	def *Definition

	HandSide HandSide

	Curses []*Curse
}

func (i *Item) Name() string {
	return i.def.Name
}

func (i *Item) Description() string {
	return fmt.Sprintf(i.def.Description, i.HandSide.String())
}

func upsertMod(mods *[]*Mod, def *Definition) *Mod {
	for _, m := range *mods {
		if def.ID == m.def.ID {
			m.Stacks++
			return m
		}
	}
	m := &Mod{
		def:    def,
		Stacks: 1,
	}
	*mods = append(*mods, m)

	return m
}

func isOneShotItem(i *Item) bool {
	switch i.def.ID {
	case Prayer, Highroll, Change_of_mind, Mimic, Trap, Lowroll, Sabotage:
		return true
	default:
		return false
	}
}

func (i *Item) RegisterMod(c *Core, phase *Phase) {
	switch {
	case isOneShotItem(i):
		// Single effect item
		var hm *HandModifiers
		switch i.HandSide {
		case RightHand:
			hm = c.right
		case LeftHand:
			hm = c.left
		}
		m := &Mod{
			def: i.def,
		}
		m.Init(c, hm, phase)
	case i.def.Hand == None:
		// Global modifiers
		switch i.def.Kind {
		case Cursed:
			upsertMod(&c.Curses, i.def).Init(c, nil, phase)
		default:
			upsertMod(&c.Bonuses, i.def).Init(c, nil, phase)
		}
	default:
		// Hand modifiers
		var hands []*HandModifiers
		switch i.HandSide {
		case RightHand:
			hands = []*HandModifiers{c.right}
		case LeftHand:
			hands = []*HandModifiers{c.left}
		case BothHand:
			hands = []*HandModifiers{c.right, c.left}
		}
		for _, h := range hands {
			switch i.def.Kind {
			case Cursed:
				upsertMod(&h.Curses, i.def).Init(c, hands[0], phase)
			default:
				upsertMod(&h.Bonuses, i.def).Init(c, hands[0], phase)
			}
		}
	}
	// Curses
	for _, ci := range i.Curses {
		(*Item)(ci).RegisterMod(c, phase)
	}
}
