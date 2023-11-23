package eval

import (
	"context"
	"errors"
)

type _let struct{}

func (_ *_let) TypeId() string {
	return "special_form.let"
}

func (_ *_let) AtomId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_let) String() string {
	return "#<syntax let>"
}

func (_ *_let) IsList() bool {
	return false
}

func (l *_let) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_let) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength < 2 {
		return nil, errors.New("malformed let")
	}

	bindings := args[0]
	body := args[1]

	bindingsArr, bindingsArrLength, err := ToArray(bindings)
	if err != nil {
		return nil, err
	}

	temporaryParamsLen := bindingsArrLength
	paramsLen := bindingsArrLength
	temporaryParams := make([]SExpression, temporaryParamsLen)
	params := make([]SExpression, paramsLen)

	for i := uint64(0); i < bindingsArrLength; i++ {
		varnameValuePair, varnameValuePairLen, varnameValuePairErr := ToArray(bindingsArr[i])
		if varnameValuePairErr != nil {
			return nil, err
		}

		if varnameValuePairLen != 2 || varnameValuePair[0].TypeId() != "symbol" {
			return nil, errors.New("malformed let")
		}

		temporaryParams = append(temporaryParams, varnameValuePair[0])
		evaluatedParams, err := Eval(ctx, varnameValuePair[1], env)
		if err != nil {
			return nil, err
		}
		params = append(params, evaluatedParams)
	}

	closure, err := NewClosure(body, temporaryParams, env, temporaryParamsLen)

	if err != nil {
		return nil, err
	}

	r, err := closure.Apply(ctx, env, params, paramsLen)

	return r, err
}

func NewLet() SExpression {
	return &_let{}
}
