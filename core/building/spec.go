package building

import (
	"github.com/Zyko0/Alapae/core/entity"
)

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

type ItemSpec struct {
	ID          int
	Rarity      Kind
	Hand        HandConstraint
	Name        string
	Description string
	MaxStacks   int

	strFn func(ctx *entity.Context) string
}

func (is *ItemSpec) String(ctx *entity.Context) string {
	return is.strFn(ctx)
}

/*
Dynamic:
- Swap [critical chance+damage/flat damage bonus/extra shot bonus/curses] from your %s with the [] from the other hand.
- Curse:
	- Gain [flat damage/crit chance] with your %s for each [accuracy/projectile speed] malus of your other hand.
	- Gain +1 flat damage with your %s for each 5% chance of NOT doing a critical strike with your other hand.
*/
