package eval

import (
	"context"
	"errors"
	"time"
)

func _subr_wait_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	if SexpressionTypeConsCell != args._sexp_type_id {
		return CreateNil(), errors.New("need arguments")
	}
	arguments := args._cell

	if !IsEmptyList(arguments._cdr) {
		return CreateNil(), errors.New("need arguments length is 1")
	}

	waitTime, err := Eval(ctx, arguments._car, env)
	if err != nil {
		return CreateNil(), err
	}

	if !waitTime.IsNumber() {
		return CreateNil(), errors.New("need 1st argument must be number but got " + waitTime._sexp_type_id.String())
	}

	durationTime := time.Millisecond * time.Duration(int(waitTime._number))
	time.Sleep(durationTime)

	return CreateNil(), nil
}

func NewWait() *Sexpression {
	return CreateSubroutine("wait", _subr_wait_Apply)
}
