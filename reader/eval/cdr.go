package eval

import (
	"context"
	"errors"
)

//this implement is lisp car function

type _cdr struct{}

func (_ *_cdr) TypeId() string {
	return "subroutine.cdr"
}

func (_ *_cdr) AtomId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_cdr) String() string {
	return "#<subr cdr>"
}

func (_ *_cdr) IsList() bool {
	return false
}

func (l *_cdr) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_cdr) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 1 > argsLength {
		return nil, errors.New("need arguments size is 1")
	}

	if args[0].TypeId() != "cons_cell" {
		return nil, errors.New("need arguments type is list")
	}

	cons := args[0].(ConsCell)

	return cons.GetCdr(), nil
}

func NewCdr() SExpression {
	return &_cdr{}
}
