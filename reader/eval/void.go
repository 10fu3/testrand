package eval

import (
	"context"
	"errors"
)

type _void struct{}

func (_ *_void) TypeId() string {
	return "special_form.void"
}

func (_ *_void) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_void) String() string {
	return "#<syntax void>"
}

func (_ *_void) IsList() bool {
	return false
}

func (l *_void) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_void) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {
	args, err := ToArray(arguments)

	if err != nil {
		return nil, err
	}

	if 1 != len(args) {
		return nil, errors.New("need arguments size is 1")
	}

	_, err = Eval(ctx, args[0], env)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func NewVoid() SExpression {
	return &_void{}
}
