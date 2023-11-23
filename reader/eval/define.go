package eval

import (
	"context"
	"errors"
)

type _define struct{}

func (_ *_define) TypeId() string {
	return "special_form.define"
}

func (_ *_define) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_define) String() string {
	return "#<syntax #define>"
}

func (_ *_define) IsList() bool {
	return false
}

func (d *_define) Equals(sexp SExpression) bool {
	return d.TypeId() == sexp.TypeId()
}

func (_ *_define) Apply(ctx context.Context, env Environment, arguments []SExpression, argsLength uint64) (SExpression, error) {

	name := arguments[0].(Symbol)

	if argsLength == 1 {
		env.Define(name, NewNil())
		return name, nil
	}

	if argsLength != 2 {
		return nil, errors.New("need less than 3 params")
	}

	initValue := arguments[1]

	evaluatedInitValue, err := Eval(ctx, initValue, env)

	if err != nil {
		return nil, err
	}
	env.Define(name, evaluatedInitValue)
	return name, nil
}

func NewDefine() SExpression {
	return &_define{}
}
