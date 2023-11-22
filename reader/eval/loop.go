package eval

import (
	"context"
	"errors"
	"fmt"
)

type _loop struct{}

func (_ *_loop) _syntax_loop_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {

	if SexpressionTypeConsCell != args.SexpressionTypeId() {
		return CreateNil(), errors.New("need arguments")
	}

	arguments := args._cell
	if arguments._cdr.SexpressionTypeId() != SexpressionTypeConsCell {
		return CreateNil(), errors.New("need arguments")
	}

	if arguments._cdr._cell._cdr.SexpressionTypeId() != SexpressionTypeConsCell {
		return CreateNil(), errors.New("need arguments for 2nd")
	}

	rawForms := arguments._cdr

	if !IsEmptyList(rawForms._cell._cdr) {
		return CreateNil(), errors.New("argument size must be 2, but got more than 2 arguments")
	}

	var evaluatedCond *Sexpression
	var err error

	for {
		evaluatedCond, err = Eval(ctx, arguments._car, env)
		if err != nil {
			return CreateNil(), err
		}

		if SexpressionTypeBool != evaluatedCond.SexpressionTypeId() {
			return CreateNil(), errors.New("need 1st argument must be bool but got " + evaluatedCond._sexp_type_id.String())
		}

		if !evaluatedCond._boolean {
			return CreateNil(), nil
		}

		_, err := Eval(ctx, rawForms._cell._car, env)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}

func NewLoop() *Sexpression {
	return CreateSpecialForm("loop", (&_loop{})._syntax_loop_Apply)
}
