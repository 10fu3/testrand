package eval

import (
	"context"
	"errors"
)

func _syntax_quote_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}
	if arrSize != 1 {
		return CreateNil(), errors.New("malformed quote")
	}
	return arr[0], nil
}

func NewQuote() *Sexpression {
	return CreateSpecialForm("quote", _syntax_quote_Apply)
}
