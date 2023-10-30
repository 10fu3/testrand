package eval

import (
	"context"
	"time"
)

type _get_now_time_micro struct{}

func (_ *_get_now_time_micro) Type() string {
	return "subroutine.get_now_time_micro"
}

func (_ *_get_now_time_micro) String() string {
	return "#<subr get_now_time_micro>"
}

func (_ *_get_now_time_micro) IsList() bool {
	return false
}

func (l *_get_now_time_micro) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_get_now_time_micro) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	return NewInt(time.Now().UnixNano()), nil
}

func NewGetNowTimeMicro() SExpression {
	return &_get_now_time_micro{}
}
