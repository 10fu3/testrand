package eval

type _and struct {
}

func (_ _and) Type() string {
	return "special_form.and"
}

func (_ _and) String() string {
	return "#<syntax #and>"
}

func (_ _and) IsList() bool {
	return false
}

func (a _and) Equals(sexp SExpression) bool {
	return a.Type() == sexp.Type()
}

func (_ _and) Apply(env Environment, args SExpression) (SExpression, error) {

	if IsEmptyList(args) {
		return NewBool(true), nil
	}

	arr, err := ToArray(args)

	if err != nil {
		return nil, err
	}

	evaluatedElm := NewConsCell(NewNil(), NewNil()).(SExpression)

	for i := 0; i < len(arr); i++ {
		evaluatedElm, err = Eval(arr[i], env)
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
