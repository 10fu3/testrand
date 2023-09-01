package eval

import (
	"context"
	"errors"
)

type _quote struct{}

func (_ *_quote) Type() string {
	return "special_form.quote"
}

func (_ *_quote) String() string {
	return "#<syntax quote>"
}

func (_ *_quote) IsList() bool {
	return false
}

func (q *_quote) Equals(sexp SExpression) bool {
	return q.Type() == sexp.Type()
}

func (_ *_quote) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)

	if err != nil {
		return nil, err
	}
	if len(arr) != 1 {
		return nil, errors.New("malformed quote")
	}
	return arr[0], nil
}

func NewQuote() SExpression {
	return &_quote{}
}
