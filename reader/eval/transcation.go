package eval

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type _transaction struct{}

func (_ *_transaction) Type() string {
	return "special_form.transaction"
}

func (_ *_transaction) String() string {
	return "#<syntax #transaction>"
}

func (_ *_transaction) IsList() bool {
	return false
}

func (t *_transaction) Equals(sexp SExpression) bool {
	return t.Type() == sexp.Type()
}

func (_ *_transaction) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {

	if "cons_cell" != args.Type() {
		return nil, errors.New("type error")
	}

	cell := args.(ConsCell)

	if !IsEmptyList(cell.GetCdr()) {
		return nil, errors.New("need less than 2 params")
	}

	ok, err := env.GetSuperGlobalEnv().Transaction(func(stm concurrency.STM) error {
		sexp, err := Eval(context.WithValue(ctx, "transaction", stm), cell.GetCar(), env)
		if err != nil {
			return err
		}
		if sexp != nil {
			fmt.Println(sexp.String())
		}
		return nil
	})

	if !ok {
		return nil, err
	}

	return NewConsCell(NewNil(), NewNil()), nil
}

func NewTransaction() SExpression {
	return &_transaction{}
}
