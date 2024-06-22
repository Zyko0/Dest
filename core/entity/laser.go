package entity

import (
	"sort"

	"github.com/Zyko0/Alapae/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MaxLaserDistance     = 512
	DefaultLaserDuration = 45
	DefaultLaserWidth    = 15
)

type Laser struct {
	ticks uint

	start    mgl64.Vec3
	end      mgl64.Vec3
	duration uint
}

func NewLaser(start, end mgl64.Vec3) *Laser {
	return &Laser{
		start:    start,
		end:      end,
		duration: DefaultLaserDuration,
	}
}

func (l *Laser) Update() {
	l.ticks = min(l.ticks+1, l.duration)
}

func (l *Laser) AppendVerticesIndices(vx []ebiten.Vertex, ix []uint16, index *int, ctx *graphics.Context) ([]ebiten.Vertex, []uint16) {
	const (
		spawnTicks = 5
	)

	t := float64(l.ticks) / float64(l.duration)
	if l.ticks > l.duration-spawnTicks {
		t = float64(l.duration-l.ticks) / spawnTicks
	}
	w := DefaultLaserWidth * 1. //t * DefaultLaserWidth
	pos := ctx.CameraPosition
	dir := l.end.Sub(l.start).Normalize()
	start := l.start.Sub(pos).Add(dir.Mul(t * 192))
	end := start.Add(dir.Mul(t * 192))
	//ls := start.Len()
	//le := end.Len()
	pv := ctx.ProjView
	s := pv.Mul4x1(start.Vec4(1))
	e := pv.Mul4x1(end.Vec4(1))
	if s.Z() < 0 || e.Z() < 0 {
		return vx, ix
	}

	const np = 0.01
	// Clip points to near plane
	if s.Z() < np {
		t := (np - s.Z()) / (e.Z() - s.Z())
		s[0] = (1-t)*s[0] + t*e[0]
		s[1] = (1-t)*s[1] + t*e[1]
		s[2] = (1-t)*s[2] + t*e[2]
		s[3] = (1-t)*s[3] + t*e[3]
	}
	if e.Z() < np {
		t := (np - e.Z()) / (s.Z() - e.Z())
		e[0] = (1-t)*e[0] + t*s[0]
		e[1] = (1-t)*e[1] + t*s[1]
		e[2] = (1-t)*e[2] + t*s[2]
		e[3] = (1-t)*e[3] + t*s[3]
	}
	for i := 0; i < 2; i++ {
		s[i] /= s.W()
		e[i] /= e.W()
	}
	sx, sy := graphics.ScreenCoordinates(s.X(), s.Y(), ctx.ScreenBounds)
	ex, ey := graphics.ScreenCoordinates(e.X(), e.Y(), ctx.ScreenBounds)

	n := mgl64.Vec2{float64(ex - sx), float64(ey - sy)}.Normalize().Mul(w)
	p0 := mgl64.Vec4{float64(sx) - n.Y(), float64(sy) + n.X(), -1, -1}
	p1 := mgl64.Vec4{float64(sx) + n.Y(), float64(sy) - n.X(), 1, -1}
	p2 := mgl64.Vec4{float64(ex) - n.Y(), float64(ey) + n.X(), -1, 1}
	p3 := mgl64.Vec4{float64(ex) + n.Y(), float64(ey) - n.X(), 1, 1}
	pts := []mgl64.Vec4{p0, p1, p2, p3}
	sort.SliceStable(pts, func(i, j int) bool {
		return mgl64.Vec2{pts[i].X(), pts[i].Y()}.LenSqr() > mgl64.Vec2{pts[j].X(), pts[j].Y()}.LenSqr()
	})
	pts[0][2], pts[0][3] = -1, -1
	pts[1][2], pts[1][3] = 1, -1
	pts[2][2], pts[2][3] = -1, 1
	pts[3][2], pts[3][3] = 1, 1
	vx, ix = graphics.AppendQuadVerticesIndices(vx, ix, *index, &graphics.QuadOpts{
		P0: pts[0],
		P1: pts[1],
		P2: pts[2],
		P3: pts[3],
		R:  2, // Laser (hardcoded)
	})
	//fmt.Println("ps", p0, p1, p2, p3)
	*index++

	return vx, ix
}

/*
func DrawRay(screen *ebiten.Image, projection elv.Projection, start, end elv.Vec3) {
	pos := projection.Position()
	start = start.Sub(pos)
	end = end.Sub(pos)
	projView := projection.ProjectionMatrix().Mul4(projection.ViewMatrix()) //.Inv()
	vs := projView.Mul4x1(start.Vec4(1))
	ve := projView.Mul4x1(end.Vec4(1))
	if vs.Z() < 0 && ve.Z() < 0 {
		return
	}
	const np = 0.01 // TODO: do not hardcode near plane
	// Clip points to near plane
	if vs.Z() < np {
		t := (np - vs.Z()) / (ve.Z() - vs.Z())
		vs = xmath.LerpVec4(vs, ve, t)
	}
	if ve.Z() < np {
		t := (np - ve.Z()) / (vs.Z() - ve.Z())
		ve = xmath.LerpVec4(ve, vs, t)
	}
	for i := 0; i < 2; i++ {
		vs[i] /= vs.W()
		ve[i] /= ve.W()
	}
	sx, sy := graphics.ScreenCoordinates(vs.X(), vs.Y(), screen.Bounds())
	ex, ey := graphics.ScreenCoordinates(ve.X(), ve.Y(), screen.Bounds())
	// TODO: clipping
	vector.StrokeLine(screen, sx, sy, ex, ey, 4., color.RGBA{255, 0, 0, 255}, true)
}*/

func (l *Laser) Dead() bool {
	return l.ticks >= l.duration
}
