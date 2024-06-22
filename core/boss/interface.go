package boss

import "github.com/Zyko0/Alapae/core/entity"

const (
	BossRadius = 2
)

type Pattern interface {
	Update(ctx *entity.Context)
	Over() bool
}

type PatternInstancier func(ctx *entity.Context) Pattern

type Sequence struct {
	index       int
	current     Pattern
	instanciers []PatternInstancier
}

func newSequence(instanciers ...PatternInstancier) *Sequence {
	return &Sequence{
		index:       0,
		instanciers: instanciers,
	}
}

func (s *Sequence) Update(ctx *entity.Context) {
	if s.current == nil || s.current.Over() {
		s.current = s.instanciers[s.index](ctx)
	}
	s.current.Update(ctx)
	if s.current.Over() {
		s.current = nil
		s.index++
	}
	if s.index >= len(s.instanciers) {
		s.index = 0
	}
}

func (s *Sequence) Over() bool {
	return false
}
