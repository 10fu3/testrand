package eval

import (
	"context"
	"time"
)

type _get_now_time_nano struct{}

func (_ *_get_now_time_nano) TypeId() string {
	return "subroutine.get_now_time_micro"
}

func (_ *_get_now_time_nano) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_get_now_time_nano) String() string {
	return "#<subr get_now_time_micro>"
}

func (_ *_get_now_time_nano) IsList() bool {
	return false
}

func (l *_get_now_time_nano) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_get_now_time_nano) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	return NewInt(time.Now().UnixNano()), nil
}

func NewGetNowTimeNano() SExpression {
	return &_get_now_time_nano{}
}
