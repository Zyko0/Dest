package core

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Camera struct {
	position  mgl64.Vec3
	direction mgl64.Vec3

	up    mgl64.Vec3
	right mgl64.Vec3

	proj mgl64.Mat4
	view mgl64.Mat4

	fov   float64
	yaw   float64
	pitch float64
}

func NewCamera(pos, front mgl64.Vec3, fov, ratio float64) *Camera {
	p := &Camera{
		position:  pos,
		direction: front,

		proj: mgl64.Perspective(
			mgl64.DegToRad(fov),
			ratio,
			0.01,
			1000, // TODO: tweak
		),

		fov:   fov,
		yaw:   0,
		pitch: 0,
	}
	p.Update()

	return p
}

func (p *Camera) FoV() float64 {
	return p.fov
}

func (p *Camera) Position() mgl64.Vec3 {
	return p.position
}

func (p *Camera) Direction() mgl64.Vec3 {
	return p.direction
}

func (p *Camera) YawPitch() (float64, float64) {
	return p.yaw, p.pitch
}

func (p *Camera) Right() mgl64.Vec3 {
	return p.right
}

func (p *Camera) ProjectionMatrix() mgl64.Mat4 {
	return p.proj
}

func (p *Camera) ViewMatrix() mgl64.Mat4 {
	return p.view
}

func (p *Camera) SetYawPitch(yaw, pitch float64) {
	p.yaw, p.pitch = yaw, pitch
}

func (p *Camera) SetPosition(position mgl64.Vec3) {
	p.position = position
}

func (p *Camera) Update() {
	const halfPi = math.Pi / 2

	p.direction = mgl64.Vec3{
		math.Cos(p.pitch) * math.Sin(p.yaw),
		math.Sin(p.pitch),
		math.Cos(p.pitch) * math.Cos(p.yaw),
	}.Normalize()
	// Right
	p.right = mgl64.Vec3{
		math.Sin(p.yaw - halfPi),
		0,
		math.Cos(p.yaw - halfPi),
	}
	// Up
	p.up = p.right.Cross(p.direction)
	p.right = p.right.Mul(-1)
	// View matrix
	pos := mgl64.Vec3{}
	target := pos.Add(p.direction)
	p.view = mgl64.LookAtV(pos, target, p.up)
}
