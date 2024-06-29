package pattern

import "github.com/Zyko0/Dest/core/entity"

type Delay struct {
	ticks    uint
	duration uint
	over     bool
}

func NewDelay(duration uint) *Delay {
	d := &Delay{
		duration: duration,
	}

	return d
}

func (d *Delay) Update(ctx *entity.Context) {
	if d.over {
		return
	}
	if d.ticks >= d.duration {
		d.over = true
		return
	}

	d.ticks++
}

func (d *Delay) Over() bool {
	return d.over
}
