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
	NextExtraCurses []*Item
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
		targetRange = 32 * 32
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

func (p *Phase) RollExisting(objects []*ItemObject, curse bool) {
	allowed := p.possibleItems()

	byRarity := [6][]int{}
	byRarity[Common] = append(byRarity[Common], commons...)
	byRarity[Uncommon] = append(byRarity[Uncommon], uncommons...)
	byRarity[Rare] = append(byRarity[Rare], rares...)
	byRarity[Epic] = append(byRarity[Epic], epics...)
	byRarity[Legendary] = append(byRarity[Legendary], legendaries...)
	byRarity[Cursed] = append(byRarity[Cursed], curses...)
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
	for i := 0; i < len(objects); i++ {
		roll := rand.Float64() * (1 - p.core.Luck)
		rarity := Common
		switch {
		case curse:
			rarity = Cursed
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
		if curse {
			objects[i].Item.Curses = append(objects[i].Item.Curses, &Item{
				def:      def,
				HandSide: hs,
			})
		} else {
			objects[i].Item.def = def
			objects[i].Item.HandSide = hs
		}
	}
}

func (p *Phase) RollNew() {
	const (
		CellSize   = float64(logic.ArenaSize) / ItemsPerLine
		ItemRadius = 0.25
		ItemHeight = -graphics.SpriteScale + ItemRadius*graphics.SpriteScale
	)

	p.Objects = make([]*ItemObject, PhaseItems)
	i := 0
	for z := 0; z < Lines; z++ {
		for x := 0; x < ItemsPerLine; x++ {
			var offx, offz float64
			if z%2 == 1 {
				offx = (CellSize / 2)
			}
			ox := float64(x)*CellSize + 2*ItemRadius*graphics.SpriteScale + offx
			oz := float64(z)*(CellSize/2) + 2*ItemRadius*graphics.SpriteScale + offz
			p.Objects[i] = &ItemObject{
				pos: mgl64.Vec3{
					ox, ItemHeight, oz,
				},
				radius: ItemRadius,
				Item:   &Item{},
			}
			i++
		}
	}
	p.RollExisting(p.Objects, false)
}

func (p *Phase) RollExtraCurses(n int) {
	items := make([]*ItemObject, n)
	for i := 0; i < n; i++ {
		items[i] = p.Objects[rand.Intn(len(p.Objects))]
	}
	p.RollExisting(items, true)
}

func (p *Phase) Pick() {
	if p.Target != nil {
		if index := slices.Index(p.Objects, p.Target); index != -1 {
			p.Target.picked = true
			p.Objects = slices.Delete(p.Objects, index, index+1)
			// Delete 3 item objects
			for i := 0; i < 3; i++ {
				index := rand.Intn(len(p.Objects))
				p.Objects[index].picked = true
				p.Objects = slices.Delete(p.Objects, index, index+1)
			}
			// Curse 10 items objects
			p.RollExtraCurses(10)
		}
		if len(p.NextExtraCurses) > 0 {
			p.Target.Item.Curses = append(p.Target.Item.Curses, p.NextExtraCurses...)
			p.NextExtraCurses = p.NextExtraCurses[:0]
		}
		p.Target = nil
	}
}

func (p *Phase) RegisterExtraCurse() {
	obj := []*ItemObject{{Item: &Item{}}}
	p.RollExisting(obj, true)
	p.NextExtraCurses = append(p.NextExtraCurses, obj[0].Item.Curses[0])
}
