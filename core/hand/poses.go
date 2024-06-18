package hand

import (
	"log"

	"github.com/Zyko0/Alapae/assets"
)

var (
	AnimationIdle        *Animation
	AnimationShootFinger *Animation
)

func init() {
	AnimationIdle = &Animation{}
	if err := AnimationIdle.Deserialize(assets.AnimIdleSrc); err != nil {
		log.Fatal("err: ", err)
	}

	AnimationShootFinger = &Animation{}
	if err := AnimationShootFinger.Deserialize(assets.AnimShootFingerSrc); err != nil {
		log.Fatal("err: ", err)
	}
}
