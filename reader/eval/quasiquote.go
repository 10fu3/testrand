package eval

import "errors"

type _quasiquote struct {
}

func (_ *_quasiquote) Type() string {
	return "special_form.quasiquote"
}

func (_ *_quasiquote) String() string {
	return "#<syntax #quasiquote>"
}

func (_ *_quasiquote) IsList() bool {
	return false
}

func (q *_quasiquote) Equals(sexp SExpression) bool {
	return q.Type() == sexp.Type()
}

func (_ *_quasiquote) Apply(env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)

	if err != nil {
		return nil, err
	}
	if len(arr) != 1 {
		return nil, errors.New("malformed quote")
	}
	return arr[0], nil
}

func NewQuasiquote() SExpression {
	return &_quasiquote{}
}
