package eval

import "context"

type _or struct {
}

func (_ _or) Type() string {
	return "special_form.and"
}

func (_ _or) String() string {
	return "#<syntax #or>"
}

func (_ _or) IsList() bool {
	return false
}

func (a _or) Equals(sexp SExpression) bool {
	return a.Type() == sexp.Type()
}

func (_ _or) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {

	if IsEmptyList(args) {
		return NewBool(false), nil
	}

	arr, err := ToArray(args)

	if err != nil {
		return nil, err
	}

	evaluatedElm := NewConsCell(NewNil(), NewNil()).(SExpression)

	for i := 0; i < len(arr); i++ {
		evaluatedElm, err = Eval(ctx, arr[i], env)
		if err != nil {
			return nil, err
		}
		if !NewBool(false).Equals(evaluatedElm) {
			return evaluatedElm, nil
		}
	}

	return evaluatedElm, nil
}

func NewOr() SExpression {
	return &_or{}
}
