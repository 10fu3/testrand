package eval

import (
	"context"
	"errors"
)

type _lambda struct{}

func (_ *_lambda) TypeId() string {
	return "special_form.lambda"
}

func (_ *_lambda) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_lambda) String() string {
	return "#<syntax lambda>"
}

func (_ *_lambda) IsList() bool {
	return false
}

func (l *_lambda) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_lambda) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 2 != argsLength {
		return nil, errors.New("need arguments size is 2")
	}

	params := args[0]
	body := args[1]

	formalsArr, formalsArrLength, err := ToArray(params)

	if err != nil {
		return nil, err
	}

	closure, err := NewClosure(body, formalsArr, env, formalsArrLength)

	if err != nil {
		return nil, err
	}

	return closure, nil
}

func NewLambda() SExpression {
	return &_lambda{}
}
