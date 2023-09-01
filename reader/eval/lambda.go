package eval

import (
	"context"
	"errors"
)

type _lambda struct{}

func (_ *_lambda) Type() string {
	return "special_form.lambda"
}

func (_ *_lambda) String() string {
	return "#<syntax lambda>"
}

func (_ *_lambda) IsList() bool {
	return false
}

func (l *_lambda) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_lambda) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {
	args, err := ToArray(arguments)
	if err != nil {
		return nil, err
	}

	if 2 != len(args) {
		return nil, errors.New("need arguments size is 2")
	}

	params := args[0]
	body := args[1]

	formalsArr, err := ToArray(params)

	if err != nil {
		return nil, err
	}

	return NewClosure(body, params, env, len(formalsArr)), nil
}

func NewLambda() SExpression {
	return &_lambda{}
}
