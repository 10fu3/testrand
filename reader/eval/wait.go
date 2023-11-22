package eval

import (
	"context"
	"errors"
	"time"
)

type _wait struct{}

func (_ *_wait) TypeId() string {
	return "special_form.wait"
}

func (_ *_wait) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_wait) String() string {
	return "#<syntax wait>"
}

func (_ *_wait) IsList() bool {
	return false
}

func (l *_wait) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_wait) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	waitTime, err := Eval(ctx, args[0], env)
	if err != nil {
		return nil, err
	}

	if waitTime.TypeId() != "number" {
		return nil, errors.New("need 1st argument must be number but got " + waitTime.TypeId())
	}
	durationTime := time.Millisecond * time.Duration(int(waitTime.(Number).GetValue()))
	time.Sleep(durationTime)

	return nil, nil
}

func NewWait() SExpression {
	return &_wait{}
}
