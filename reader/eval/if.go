package eval

import (
	"context"
	"errors"
	"fmt"
)

type _if struct {
}

func NewIf() SExpression {
	return &_if{}
}

func (_ *_if) TypeId() string {
	return "special_form.if"
}

func (_ *_if) AtomId() AtomType {
	return AtomTypeSpecialForm
}

func (_ *_if) String() string {
	return "#<syntax #if>"
}

func (_ *_if) IsList() bool {
	return false
}

func (i *_if) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_if) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength <= 1 || argsLength >= 4 {
		return nil, errors.New(fmt.Sprintf("too many argument: %d", len(args)))
	}
	argsIndex := 0

	statement := args[argsIndex]
	argsIndex++
	onTrue := args[argsIndex]

	evaluated, err := Eval(ctx, statement, env)
	if err != nil {
		return Eval(ctx, statement, env)
	}

	if evaluated.Equals(NewBool(false)) {
		if argsLength == 2 {
			return NewNil(), nil
		}
		argsIndex++
		onFalse := args[argsIndex]
		return Eval(ctx, onFalse, env)
	}
	return Eval(ctx, onTrue, env)
}
