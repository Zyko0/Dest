package hand

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Zyko0/Ebiary/asset"
	"github.com/go-gl/mathgl/mgl64"
)

type Orientation struct {
	XY float64
	Z  float64
}

func (o0 Orientation) Lerp(o1 Orientation, t float64) Orientation {
	return Orientation{
		XY: (1-t)*o0.XY + t*o1.XY,
		Z:  (1-t)*o0.Z + t*o1.Z,
	}
}

func (o *Orientation) Serialize() []byte {
	return []byte(fmt.Sprintf("%.8f,%.8f", o.XY, o.Z))
}

func (o *Orientation) Deserialize(data []byte) error {
	n, err := fmt.Sscanf(string(data), "%f,%f", &o.XY, &o.Z)
	if err != nil {
		return err
	}
	if n != 2 {
		return fmt.Errorf("err: can't deserialize orientation: (%d != 2) fields", n)
	}

	return nil
}

type FingerStep [3]Orientation

func (f0 *FingerStep) Lerp(f1 *FingerStep, t float64) *FingerStep {
	f := &FingerStep{}
	f[0] = f0[0].Lerp(f1[0], t)
	f[1] = f0[1].Lerp(f1[1], t)
	f[2] = f0[2].Lerp(f1[2], t)

	return f
}

func (f *FingerStep) Serialize() []byte {
	return []byte(fmt.Sprintf("%s;%s;%s",
		f[0].Serialize(), f[1].Serialize(), f[2].Serialize(),
	))
}

func (f *FingerStep) Deserialize(data []byte) error {
	parts := strings.Split(string(data), ";")
	if len(parts) != 3 {
		return fmt.Errorf("err: can't deserialize fingerstep: (%d != 3) bones", len(parts))
	}
	for i := range *f {
		if err := f[i].Deserialize([]byte(parts[i])); err != nil {
			return err
		}
	}

	return nil
}

type HandStep struct {
	Duration uint
	Rotation mgl64.Vec3
	Fingers  [5]*FingerStep
}

func (h0 *HandStep) Lerp(h1 *HandStep, t float64) *HandStep {
	h := &HandStep{}
	h.Rotation[0] = (1-t)*h0.Rotation[0] + t*h1.Rotation[0]
	h.Rotation[1] = (1-t)*h0.Rotation[1] + t*h1.Rotation[1]
	h.Rotation[2] = (1-t)*h0.Rotation[2] + t*h1.Rotation[2]
	for i := range h.Fingers {
		h.Fingers[i] = h0.Fingers[i].Lerp(h1.Fingers[i], t)
	}
	return h
}

func (h *HandStep) Serialize() []byte {
	var fingers string
	for i := range h.Fingers {
		fingers += string(h.Fingers[i].Serialize()) + "\n"
	}

	return []byte(fmt.Sprintf("%d,%.8f,%.8f,%.8f,{\n%s}",
		h.Duration, h.Rotation[0], h.Rotation[1], h.Rotation[2], fingers,
	))
}

func (h *HandStep) Deserialize(data []byte) error {
	n, err := fmt.Sscanf(string(data), "%d,%f,%f,%f,{\n",
		&h.Duration, &h.Rotation[0], &h.Rotation[1], &h.Rotation[2],
	)
	if err != nil {
		return err
	}
	if n != 4 {
		return fmt.Errorf("err: can't deserialize handstep: (%d != 4) fields", n)
	}
	nl := strings.IndexRune(string(data), '\n')
	if nl == -1 {
		return errors.New("err: EOF but expected finger information")
	}
	fingers := strings.Trim(string(data[nl:]), "\n}")
	parts := strings.Split(fingers, "\n")
	if len(parts) != 5 {
		return fmt.Errorf("err: can't deserialize handstep: (%d != 5) fingers", len(parts))
	}
	for i := range h.Fingers {
		h.Fingers[i] = &FingerStep{}
		if err := h.Fingers[i].Deserialize([]byte(parts[i])); err != nil {
			return err
		}
	}

	return nil
}

type Animation struct {
	Loop  int
	Init  int
	Steps []*HandStep
}

func (a *Animation) Serialize() []byte {
	var steps string

	for _, s := range a.Steps {
		steps += string(s.Serialize()) + "/\n"
	}
	return []byte(fmt.Sprintf("loop:%d,init:%d\n%s", a.Loop, a.Init, steps))
}

func (a *Animation) Deserialize(data []byte) error {
	n, err := fmt.Sscanf(string(data), "loop:%d,init:%d\n", &a.Loop, &a.Init)
	if err != nil {
		return err
	}
	if n != 2 {
		return errors.New("err: can't deserialize animation: expected loop,init")
	}
	nl := strings.IndexRune(string(data), '\n')
	if nl == -1 {
		return errors.New("err: EOF but expected finger information")
	}
	steps := strings.Trim(string(data)[nl:], "\n")
	parts := strings.Split(steps, "/")
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		s := &HandStep{}
		if err := s.Deserialize([]byte(strings.Trim(parts[i], "\n"))); err != nil {
			return err
		}
		a.Steps = append(a.Steps, s)
	}

	return nil
}

type AnimationInstance struct {
	reverted  bool
	ticks     uint
	index     int
	prev      *HandStep
	loops     int
	animation *Animation
}

func (a *Animation) NewInstance(hand *Hand, reverted bool) *AnimationInstance {
	prev := hand.AsStep()
	prev.Duration = uint(a.Init)
	return &AnimationInstance{
		ticks:     0,
		index:     0,
		prev:      prev,
		animation: a,
	}
}

var TestAnimation *asset.LiveAsset[*Animation]

func init() {
	var err error

	TestAnimation, err = asset.NewLiveAssetFunc("assets/animations/run.pose", func(b []byte) (*Animation, error) {
		a := &Animation{}
		if err := a.Deserialize(b); err != nil {
			return nil, err
		}
		return a, nil
	})
	if err != nil {
		log.Fatal("animation err: ", err)
	}
}

func (a *AnimationInstance) Update(hand *Hand) {
	if a == nil {
		return
	}
	/*a.animation = TestAnimation.Value() // TODO: tmp
	if err := TestAnimation.Error(); err != nil {
		fmt.Println("err:", err)
	}*/
	var s *HandStep
	// Lerp
	idx := a.index
	if a.reverted {
		idx = len(a.animation.Steps) - idx
	}
	t := float64(a.ticks) / float64(a.prev.Duration)
	s = a.prev.Lerp(a.animation.Steps[idx], t)
	// Update hand
	hand.Rotation = s.Rotation
	for i := range hand.Fingers {
		for j := range hand.Fingers[i].Bones {
			hand.Fingers[i].Bones[j].XY = s.Fingers[i][j].XY
			hand.Fingers[i].Bones[j].Z = s.Fingers[i][j].Z
		}
	}
	a.ticks++
	// Next step
	if a.ticks >= a.prev.Duration {
		a.prev = a.animation.Steps[idx]
		a.ticks = 0
		a.index++
	}
	// Reset to loop index
	if a.index >= len(a.animation.Steps) {
		a.index = a.animation.Loop
		a.loops++
	}
}

func (a *AnimationInstance) Loops() int {
	return a.loops
}
