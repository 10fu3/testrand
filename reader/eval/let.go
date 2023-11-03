package eval

import (
	"context"
	"errors"
)

type _let struct{}

func (_ *_let) Type() string {
	return "special_form.let"
}

func (_ *_let) String() string {
	return "#<syntax let>"
}

func (_ *_let) IsList() bool {
	return false
}

func (l *_let) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_let) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)

	if err != nil {
		return nil, err
	}
	if len(arr) < 2 {
		return nil, errors.New("malformed let")
	}

	bindings := arr[0]
	body := arr[1]

	bindingsArr, err := ToArray(bindings.(ConsCell).GetCar())
	if err != nil {
		return nil, err
	}

	if len(bindingsArr)%2 != 0 {
		return nil, errors.New("malformed let")
	}

	temporaryParams := make([]SExpression, 0)
	params := make([]SExpression, 0)

	for i := 0; i < len(bindingsArr); i += 2 {
		temporaryParams = append(temporaryParams, bindingsArr[i])
		evaluatedParams, err := Eval(ctx, bindingsArr[i+1], env)
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
