//go:build !devx

package assets

import (
	"log"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed shaders/arena.kage
	shaderArenaSrc []byte
	shaderArena    *ebiten.Shader

	//go:embed shaders/hand.kage
	shaderHandsSrc []byte
	shaderHands    *ebiten.Shader

	//go:embed shaders/marker.kage
	shaderMarkerSrc []byte
	shaderMarker    *ebiten.Shader

	//go:embed shaders/entity.kage
	shaderEntitySrc []byte
	shaderEntity    *ebiten.Shader
)

func init() {
	var err error

	shaderArena, err = ebiten.NewShader(shaderArenaSrc)
	if err != nil {
		log.Fatal("shader: ", err)
	}

	shaderHands, err = ebiten.NewShader(shaderHandsSrc)
	if err != nil {
		log.Fatal("shader: ", err)
	}

	shaderMarker, err = ebiten.NewShader(shaderMarkerSrc)
	if err != nil {
		log.Fatal("shader: ", err)
	}

	shaderEntity, err = ebiten.NewShader(shaderEntitySrc)
	if err != nil {
		log.Fatal("shader: ", err)
	}
}

func ShaderArena() *ebiten.Shader {
	return shaderArena
}

func ShaderHands() *ebiten.Shader {
	return shaderHands
}

func ShaderMarker() *ebiten.Shader {
	return shaderMarker
}

func ShaderEntity() *ebiten.Shader {
	return shaderEntity
}
