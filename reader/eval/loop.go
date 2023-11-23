package eval

import (
	"context"
	"errors"
	"fmt"
)

type _loop struct{}

func (_ *_loop) TypeId() string {
	return "special_form.loop"
}

func (_ *_loop) AtomId() AtomType {
	return AtomTypeSpecialForm
}

func (_ *_loop) String() string {
	return "#<syntax loop>"
}

func (_ *_loop) IsList() bool {
	return false
}

func (l *_loop) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_loop) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New(fmt.Sprintf("malformed loop: %d", len(args)))
	}

	var evaluatedCond SExpression
	var err error

	for {
		evaluatedCond, err = Eval(ctx, args[0], env)
		if err != nil {
			return nil, err
		}

		if evaluatedCond.TypeId() != "bool" {
			return nil, errors.New("need 1st argument must be bool but got " + evaluatedCond.TypeId())
		}

		if !evaluatedCond.(Bool).GetValue() {
			return nil, nil
		}

		_, err := Eval(ctx, args[1], env)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}

func NewLoop() SExpression {
	return &_loop{}
}
