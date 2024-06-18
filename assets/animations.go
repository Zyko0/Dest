package assets

import (
	_ "embed"
)

var (
	//go:embed animations/idle.pose
	AnimIdleSrc []byte

	//go:embed animations/shoot_finger.pose
	AnimShootFingerSrc []byte
)
