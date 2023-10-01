package eval

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"
	"strings"
)

type _global_get struct {
}

func (_ *_global_get) Type() string {
	return "special_form.global_set"
}

func (_ *_global_get) String() string {
	return "#<syntax global_set>"
}

func (_ *_global_get) IsList() bool {
	return false
}

func (l *_global_get) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_global_get) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	if "cons_cell" != args.Type() {
		return nil, errors.New("type error")
	}

	cell := args.(ConsCell)

	name := cell.GetCar().(Symbol)

	var err error

	if err != nil {
		return nil, err
	}

	result, err := func() (SExpression, error) {
		var err error
		var result SExpression
		if ctx.Value("transaction") != nil {
			transaction := ctx.Value("transaction").(concurrency.STM)
			var r = transaction.Get(fmt.Sprintf("/env/%s/%s", env.GetParentId(), name.String()))
			input := strings.NewReader(fmt.Sprintf("%s\n", r))
			reader := New(bufio.NewReader(input))
			result, err = reader.Read()
		} else {
			r, err := env.GetSuperGlobalEnv().Get(fmt.Sprintf("/env/%s/%s", env.GetParentId(), name.String()))
			if err != nil {
				return nil, err
			}
			input := strings.NewReader(fmt.Sprintf("%s\n", r))
			reader := New(bufio.NewReader(input))
			result, err = reader.Read()
		}
		return result, err
	}()

	return result, err
}

func NewGlobalGet() SExpression {
	return &_global_get{}
}
