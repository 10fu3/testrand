package eval

import (
	"context"
	"errors"
)

type _quote struct{}

func (_ *_quote) TypeId() string {
	return "special_form.quote"
}

func (_ *_quote) AtomId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_quote) String() string {
	return "#<syntax quote>"
}

func (_ *_quote) IsList() bool {
	return false
}

func (q *_quote) Equals(sexp SExpression) bool {
	return q.TypeId() == sexp.TypeId()
}

func (_ *_quote) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	if argsLength != 1 {
		return nil, errors.New("malformed quote")
	}
	return args[0], nil
}

func NewQuote() SExpression {
	return &_quote{}
}
