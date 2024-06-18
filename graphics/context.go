package graphics

import "github.com/go-gl/mathgl/mgl64"

type Context struct {
	CameraPosition mgl64.Vec3
	CameraRight    mgl64.Vec3
	CameraUp       mgl64.Vec3
	ProjView       mgl64.Mat4
	ViewInv        mgl64.Mat4
}
