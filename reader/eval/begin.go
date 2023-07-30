package eval

type _begin struct{}

func (_ *_begin) Type() string {
	return "special_form.begin"
}

func (_ *_begin) String() string {
	return "#<syntax #begin>"
}

func (_ *_begin) IsList() bool {
	return false
}

func (b *_begin) Equals(sexp SExpression) bool {
	return b.Type() == sexp.Type()
}

func (_ *_begin) Apply(env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)

	if err != nil {
		return nil, err
	}

	last := NewNil()

	for _, sexp := range arr {
		last, err = Eval(sexp, env)
		if err != nil {
			return nil, err
		}
	}
	return last, err
}

func NewBegin() SExpression {
	return &_begin{}
}
