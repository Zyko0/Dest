package building

import (
	"math/rand"
	"slices"

	"github.com/Zyko0/Alapae/core/aoe"
	"github.com/Zyko0/Alapae/core/entity"
	"github.com/Zyko0/Alapae/core/hand"
	"github.com/Zyko0/Alapae/graphics"
	"github.com/Zyko0/Alapae/logic"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	ItemsPerLine = 12
	Lines        = 20
	PhaseItems   = Lines * ItemsPerLine
)

type Phase struct {
	core *Core

	Target          *ItemObject
	NextExtraCurses []*Curse
	Objects         []*ItemObject
}

func NewPhase(c *Core) *Phase {
	return &Phase{
		core: c,
	}
}

func (p *Phase) MarkerShape() aoe.Shape {
	if p.Target == nil {
		return nil
	}

	return &aoe.CircleBorder{
		X:      float32(p.Target.pos.X()),
		Y:      float32(p.Target.pos.Z()),
		Radius: float32(2 * p.Target.radius * graphics.SpriteScale),
	}
}

func (p *Phase) Update(ctx *entity.Context) {
	const (
		targetRange = 10 * 10 * 10 * 10
	)
	var bestDistSq = 99999.

	p.Target = nil
	for _, o := range p.Objects {
		o.Update(ctx)
		if o.targeted {
			distSq := o.pos.Sub(ctx.PlayerPosition).LenSqr()
			if distSq < targetRange && distSq < bestDistSq {
				bestDistSq = distSq
				p.Target = o
			}
		}
	}
}

func (p *Phase) possibleItems() [ItemMax][2]int {
	var allowed [ItemMax][2]int

	for id := range allowed {
		switch impls[id].Hand {
		case SingleHand:
			allowed[id][hand.Right] = p.core.right.allowedStacks(id)
			allowed[id][hand.Left] = p.core.left.allowedStacks(id)
		case BothHands, None:
			allowed[id][0] = p.core.allowedStacks(id)
		}
	}

	return allowed
}

func (p *Phase) AppendEntities(entities []entity.Entity) []entity.Entity {
	for _, o := range p.Objects {
		entities = append(entities, o)
	}
	return entities
}

var (
	itemChances = [5]float64{
		Common:    1,
		Uncommon:  0.9,
		Rare:      0.5,
		Epic:      0.2,
		Legendary: 0.05,
	}
)

func (p *Phase) RollExisting() {
	allowed := p.possibleItems()

	byRarity := [5][]int{}
	byRarity[Common] = append(byRarity[Common], commons...)
	byRarity[Uncommon] = append(byRarity[Uncommon], uncommons...)
	byRarity[Rare] = append(byRarity[Rare], rares...)
	byRarity[Epic] = append(byRarity[Epic], epics...)
	byRarity[Legendary] = append(byRarity[Legendary], legendaries...)
	for i, items := range byRarity {
		n := 0
		for _, id := range items {
			if allowed[id][0] > 0 || allowed[id][1] > 0 {
				items[n] = id
				n++
			}
		}
		byRarity[i] = items[:n]
	}
	var items []int
	for i := 0; i < len(p.Objects); i++ {
		roll := rand.Float64() * (1 - p.core.Luck)
		rarity := Common
		switch {
		case roll < itemChances[Legendary]:
			rarity = Legendary
		case roll < itemChances[Epic]:
			rarity = Epic
		case roll < itemChances[Rare]:
			rarity = Rare
		case roll < itemChances[Uncommon]:
			rarity = Uncommon
		default:
			rarity = Common
		}
		items = byRarity[rarity]
		index := rand.Intn(len(items))
		id := items[index]
		def := impls[id]
		var hs HandSide
		switch def.Hand {
		case SingleHand:
			switch {
			case allowed[id][0] > 0 && allowed[id][1] > 0:
				hs = HandSide(rand.Intn(2))
			case allowed[id][hand.Right] > 0:
				hs = RightHand
			case allowed[id][hand.Left] > 0:
				hs = LeftHand
			}
			allowed[id][hs] = max(allowed[id][hs]-1, 0)
		case BothHands:
			hs = BothHand
			allowed[id][0] = max(allowed[id][0]-1, 0)
		case None:
			hs = NoHand
			allowed[id][0] = max(allowed[id][0]-1, 0)
		}
		if allowed[id][0] == 0 && allowed[id][1] == 0 {
			byRarity[rarity] = slices.Delete(items, index, index+1)
		}
		p.Objects[i].Item.def = def
		p.Objects[i].Item.HandSide = hs
	}
}

func (p *Phase) RollNew() {
	const (
		CellSize   = float64(logic.ArenaSize) / (ItemsPerLine)
		ItemRadius = 0.25
		ItemHeight = -graphics.SpriteScale + ItemRadius*graphics.SpriteScale
	)

	p.Objects = make([]*ItemObject, PhaseItems)
	for i := range p.Objects {
		p.Objects[i] = &ItemObject{
			pos: mgl64.Vec3{
				float64((i%(ItemsPerLine)))*CellSize + ItemRadius*graphics.SpriteScale,
				ItemHeight,
				float64((i/(ItemsPerLine*2)))*CellSize + ItemRadius*graphics.SpriteScale,
			},
			radius: ItemRadius,
			Item:   &Item{},
		}
	}
	p.RollExisting()
}

func (p *Phase) RollExtraCurses(n int) {

}

func (p *Phase) RegisterExtraCurse() {
	//p.ExtraCurses = append(p.ExtraCurses, )
	// TODO:
}
