package eval

import "context"

type _and struct {
}

func (_ _and) TypeId() string {
	return "special_form.and"
}

func (_ _and) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ _and) String() string {
	return "#<syntax #and>"
}

func (_ _and) IsList() bool {
	return false
}

func (a _and) Equals(sexp SExpression) bool {
	return a.TypeId() == sexp.TypeId()
}

func (_ _and) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {

	if IsEmptyList(args) {
		return NewBool(true), nil
	}

	arr, err := ToArray(args)

	if err != nil {
		return nil, err
	}

	evaluatedElm := NewNil()

	for i := 0; i < len(arr); i++ {
		evaluatedElm, err = Eval(ctx, arr[i], env)
		if err != nil {
			return nil, err
		}
		if NewBool(false).Equals(evaluatedElm) {
			return evaluatedElm, nil
		}
	}

	return evaluatedElm, nil
}

func NewAnd() SExpression {
	return &_and{}
}
