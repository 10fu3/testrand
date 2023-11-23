package eval

import "context"

type _global_clear_all struct{}

func (_ *_global_clear_all) TypeId() string {
	return "subroutine.global_clear_all"
}

func (_ *_global_clear_all) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_global_clear_all) String() string {
	return "#<subr global_clear_all>"
}

func (_ *_global_clear_all) IsList() bool {
	return false
}

func (l *_global_clear_all) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_global_clear_all) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	if err := env.GetSuperGlobalEnv().ClearAll(); err != nil {
		return nil, err
	}
	return NewNil(), nil
}

func NewGlobalClearAll() SExpression {
	return &_global_clear_all{}
}
