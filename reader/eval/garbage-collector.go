package eval

import (
	"context"
	"runtime"
)

type _force_gc struct{}

func (s *_force_gc) TypeId() string {
	return "subroutine.force_gc"
}

func (s *_force_gc) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (s *_force_gc) String() string {
	return "#<subr force_gc>"
}

func (s *_force_gc) IsList() bool {
	return false
}

func (s *_force_gc) Equals(sexp SExpression) bool {
	return s.TypeId() == sexp.TypeId()
}

func (s *_force_gc) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	runtime.GC()
	return NewBool(true), nil
}

func NewForceGC() SExpression {
	return &_force_gc{}
}
