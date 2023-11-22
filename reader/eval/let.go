package eval

import (
	"context"
	"errors"
)

func _syntax_let_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {

	if SexpressionTypeConsCell != args.SexpressionTypeId() {
		return CreateNil(), errors.New("malformed let")
	}

	arr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}
	if arrSize < 2 {
		return CreateNil(), errors.New("malformed let")
	}

	bindings := arr[0]
	body := arr[1]

	bindingsArr, bindingsArrSize, err := ToArray(bindings)
	if err != nil {
		return CreateNil(), err
	}

	temporaryParams := make([]*Sexpression, 0)
	params := make([]*Sexpression, 0)

	for i := 0; i < len(bindingsArr); i++ {
		varnameValuePair, varnameValuePairLen, toArrayErr := ToArray(bindingsArr[i])
		if toArrayErr != nil {
			return CreateNil(), toArrayErr
		}

		if varnameValuePairLen != 2 || !varnameValuePair[0].IsSymbol() {
			return CreateNil(), errors.New("malformed let")
		}

		temporaryParams = append(temporaryParams, varnameValuePair[0])
		evaluatedParams, evaluatedErr := Eval(ctx, varnameValuePair[1], env)
		if evaluatedErr != nil {
			return CreateNil(), err
		}
		params = append(params, evaluatedParams)
	}

	temporaryParamsSize := bindingsArrSize

	paramsForConsCell := ToConsCell(params)

	closure, generateClosureErr := CreateClosure(body, temporaryParams, env, temporaryParamsSize)

	if generateClosureErr != nil {
		return CreateNil(), generateClosureErr
	}

	r, err := closure._applyFunc(closure, ctx, env, paramsForConsCell)

	return r, err
}

func NewLet() *Sexpression {
	return CreateSpecialForm("let", _syntax_let_Apply)
}
