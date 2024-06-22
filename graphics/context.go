package graphics

import (
	"image"

	"github.com/go-gl/mathgl/mgl64"
)

const (
	SpriteScale = 10
)

type Context struct {
	ScreenBounds   image.Rectangle
	CameraPosition mgl64.Vec3
	CameraRight    mgl64.Vec3
	CameraUp       mgl64.Vec3
	ProjView       mgl64.Mat4
	ViewInv        mgl64.Mat4
}
