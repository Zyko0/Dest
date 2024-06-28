//go:build devx

package assets

import (
	"log"

	"github.com/Zyko0/Ebiary/asset"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	shaderArena  *asset.LiveAsset[*ebiten.Shader]
	shaderHands  *asset.LiveAsset[*ebiten.Shader]
	shaderMarker *asset.LiveAsset[*ebiten.Shader]
	shaderEntity *asset.LiveAsset[*ebiten.Shader]
)

func init() {
	var err error

	shaderArena, err = asset.NewLiveAsset[*ebiten.Shader]("assets/shaders/arena.kage")
	if err != nil {
		log.Fatal("shader: ", err)
	}

	shaderHands, err = asset.NewLiveAsset[*ebiten.Shader]("assets/shaders/hand.kage")
	if err != nil {
		log.Fatal("shader: ", err)
	}

	shaderMarker, err = asset.NewLiveAsset[*ebiten.Shader]("assets/shaders/marker.kage")
	if err != nil {
		log.Fatal("shader: ", err)
	}

	shaderEntity, err = asset.NewLiveAsset[*ebiten.Shader]("assets/shaders/entity.kage")
	if err != nil {
		log.Fatal("shader: ", err)
	}
}

func ShaderArena() *ebiten.Shader {
	return shaderArena.Value()
}

func ShaderHands() *ebiten.Shader {
	return shaderHands.Value()
}

func ShaderMarker() *ebiten.Shader {
	return shaderMarker.Value()
}

func ShaderEntity() *ebiten.Shader {
	return shaderEntity.Value()
}
