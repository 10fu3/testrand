package eval

type _heavy struct{}

func (_ *_heavy) Type() string {
	return "special_form.heavy"
}

func (_ *_heavy) String() string {
	return "#<syntax heavy>"
}

func (_ *_heavy) IsList() bool {
	return false
}

func (h *_heavy) Equals(sexp SExpression) bool {
	return h.Type() == sexp.Type()
}

func (_ *_heavy) Apply(env Environment, arguments SExpression) (SExpression, error) {
	args, err := ToArray(arguments)

	if err != nil {
		return nil, err
	}

	if 1 == len(args) {

	}
	return nil, err
}
