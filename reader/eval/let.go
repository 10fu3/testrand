package eval

import (
	"context"
	"errors"
)

type _let struct{}

func (_ *_let) TypeId() string {
	return "special_form.let"
}

func (_ *_let) SExpressionTypeId() SExpressionType {
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

func (_ *_let) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {

	if args.TypeId() != "cons_cell" {
		return nil, errors.New("malformed let")
	}

	arr, err := ToArray(args.(ConsCell))

	if err != nil {
		return nil, err
	}
	if len(arr) < 2 {
		return nil, errors.New("malformed let")
	}

	bindings := arr[0]
	body := arr[1]

	bindingsArr, err := ToArray(bindings)
	if err != nil {
		return nil, err
	}

	temporaryParams := make([]SExpression, 0)
	params := make([]SExpression, 0)

	for i := 0; i < len(bindingsArr); i++ {
		varnameValuePair, err := ToArray(bindingsArr[i])
		if err != nil {
			return nil, err
		}

		if len(varnameValuePair) != 2 || varnameValuePair[0].TypeId() != "symbol" {
			return nil, errors.New("malformed let")
		}

		temporaryParams = append(temporaryParams, varnameValuePair[0])
		evaluatedParams, err := Eval(ctx, varnameValuePair[1], env)
		if err != nil {
			return nil, err
		}
		params = append(params, evaluatedParams)
	}

	temporaryParamsForConsCell := ToConsCell(temporaryParams)
	paramsForConsCell := ToConsCell(params)

	r, err := (NewClosure(body, temporaryParamsForConsCell, env, len(temporaryParams))).Apply(ctx, env, paramsForConsCell)

	return r, err
}

func NewLet() SExpression {
	return &_let{}
}
