package eval

import (
	"context"
	"errors"
)

type _set struct{}

func (_ *_set) TypeId() string {
	return "special_form.set"
}

func (_ *_set) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_set) String() string {
	return "#<syntax #set>"
}

func (_ *_set) IsList() bool {
	return false
}

func (s *_set) Equals(sexp SExpression) bool {
	return s.TypeId() == sexp.TypeId()
}

func (_ *_set) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("malformed set")
	}

	name := args[0].(Symbol)

	evaluatedInitValue, err := Eval(ctx, args[1], env)

	if err != nil {
		return nil, err
	}

	err = env.Set(name, evaluatedInitValue)
	if err != nil {
		return nil, err
	}
	return name, nil
}

func NewSet() SExpression {
	return &_set{}
}
