package eval

import (
	"context"
	"errors"
)

func _subr_is_equals_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	if SexpressionTypeConsCell != arguments._sexp_type_id {
		return CreateNil(), errors.New("type error")
	}
	argCell := arguments._cell

	first := argCell._car

	if SexpressionTypeConsCell != argCell._cdr.SexpressionTypeId() {
		return CreateNil(), errors.New("type error")
	}

	second := argCell._cdr._cell

	if !IsEmptyList(second._cdr) {
		return CreateNil(), errors.New("argument size error")
	}

	return CreateBool(first.Equals(second._car)), nil
}

func NewIsEq() *Sexpression {
	return CreateSubroutine("eq?", _subr_is_equals_Apply)
}
