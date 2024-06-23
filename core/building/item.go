package building

import (
	"image"

	"github.com/Zyko0/Alapae/core/hand"
	"github.com/go-gl/mathgl/mgl64"
)

type HandSide byte

const (
	RightHand HandSide = iota
	LeftHand
	BothHand
)

type Item struct {
	spec *ItemSpec

	Position mgl64.Vec3

	SpriteRect image.Rectangle
	HandSide   hand.Side
	Stacks     int

	Curses []*Curse
}
