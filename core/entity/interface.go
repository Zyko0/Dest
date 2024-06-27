package entity

import (
	"github.com/Zyko0/Alapae/core/aoe"
	"github.com/Zyko0/Alapae/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

type Team byte

const (
	TeamAlly Team = iota
	TeamEnemy
	TeamNone
)

type Entity interface {
	Update(ctx *Context)
	AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, ctx *graphics.Context) ([]ebiten.Vertex, []uint16)

	Team() Team
	Position() mgl64.Vec3
	Radius() float64
	Dead() bool
}

type Stance byte

const (
	StanceIdle Stance = iota
	StanceHostile
)

type Boss interface {
	Entity

	SetPosition(pos mgl64.Vec3)
	SetStance(stance Stance)

	Image() *ebiten.Image
	MarkerShape() *aoe.CircleBorder
}

type Context struct {
	CameraRight     mgl64.Vec3
	CameraUp        mgl64.Vec3
	PlayerPosition  mgl64.Vec3
	PlayerDirection mgl64.Vec3
	Boss            Boss

	Entities []Entity
	Markers  []*aoe.Marker
}
