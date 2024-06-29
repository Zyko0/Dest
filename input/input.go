package input

import (
	"runtime"

	"github.com/Zyko0/Alapae/logic"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	inputGiven bool
	cursorInit bool

	lcx, lcy int
)

func SetLastCursor(cx, cy int) {
	lcx, lcy = cx, cy
}

func EnsureCursorCaptured() bool {
	if ebiten.CursorMode() == ebiten.CursorModeCaptured {
		return true
	}

	inputGiven = inputGiven || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if runtime.GOOS != "js" || inputGiven {
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		inputGiven = false
	}

	return false
}

func ProcessMouseMovement() (float64, float64) {
	var xoff, yoff int
	cx, cy := ebiten.CursorPosition()
	// Note: ebitengine hack, with mouse captured initial cursor position is wrong
	if !cursorInit {
		if cx != 0 && cy != 0 {
			cursorInit = true
		}
	} else if cx != lcx || cy != lcy {
		xoff, yoff = cx-lcx, lcy-cy
	}
	lcx, lcy = cx, cy

	sens := logic.MouseSensitivity/1000
	x := float64(xoff) * sens
	y := float64(yoff) * sens

	return x, y
}

func ProcessKeyboard(pos, dir, right mgl64.Vec3, ms float64) mgl64.Vec3 {
	const diagMult = 0.707106

	if (ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS)) &&
		(ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyA)) {
		ms *= diagMult
	}
	d := mgl64.Vec2{dir.X(), dir.Z()}.Normalize().Mul(ms)
	r := mgl64.Vec2{right.X(), right.Z()}.Normalize().Mul(ms)
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		pos = pos.Add(mgl64.Vec3{d.X(), 0, d.Y()})
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		pos = pos.Sub(mgl64.Vec3{d.X(), 0, d.Y()})
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		pos = pos.Add(mgl64.Vec3{r.X(), 0, r.Y()})
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		pos = pos.Sub(mgl64.Vec3{r.X(), 0, r.Y()})
	}

	return pos
}
