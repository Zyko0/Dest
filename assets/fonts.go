package assets

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"

	_ "embed"
)

var (
	//go:embed fonts/ChakraPetch-SemiBold.ttf
	ttf []byte

	FontSource *text.GoTextFaceSource
)

func init() {
	var err error

	FontSource, err = text.NewGoTextFaceSource(bytes.NewReader(ttf))
	if err != nil {
		log.Fatal("err: can't create face source: ", err)
	}
}
