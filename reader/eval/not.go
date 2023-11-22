package eval

import (
	"context"
	"errors"
)

func _subr_is_not_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {

	if SexpressionTypeConsCell != arguments.SexpressionTypeId() {
		return CreateNil(), errors.New("need arguments")
	}

	argCell := arguments._cell

	first := argCell._car

	if SexpressionTypeBool != first.SexpressionTypeId() {
		return first, nil
	}

	return CreateBool(!first._boolean), nil
}

func NewIsNot() *Sexpression {
	return CreateSubroutine("not", _subr_is_not_Apply)
}
