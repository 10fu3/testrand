package eval

import (
	"bufio"
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"
	"strings"
)

type _global_get struct {
}

func (_ *_global_get) TypeId() string {
	return "special_form.global_get"
}

func (_ *_global_get) AtomId() AtomType {
	return AtomTypeSpecialForm
}

func (_ *_global_get) String() string {
	return "#<syntax global_get>"
}

func (_ *_global_get) IsList() bool {
	return false
}

func (l *_global_get) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_global_get) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	name := args[0].(Symbol)

	defaultArg := (func() SExpression {
		if argsLength == 1 {
			return NewConsCell(NewNil(), NewNil())
		}
		return args[1]
	})()

	var err error

	if err != nil {
		return nil, err
	}

	var result SExpression
	if ctx.Value("transaction") != nil {
		transaction := ctx.Value("transaction").(concurrency.STM)
		existKey := transaction.Rev(fmt.Sprintf("/env/%s/%s", env.GetParentId(), name.String()))

		if defaultArg.TypeId() != "nil" && existKey == 0 {
			return defaultArg, nil
		}

		var r = transaction.Get(fmt.Sprintf("/env/%s/%s", env.GetParentId(), name.String()))
		input := strings.NewReader(fmt.Sprintf("%s\n", r))
		reader := New(bufio.NewReader(input))
		result, err = reader.Read()
		reader = nil

	} else {
		r, err := env.GetSuperGlobalEnv().Get(fmt.Sprintf("/env/%s/%s", env.GetParentId(), name.String()))
		if err != nil {
			return nil, err
		}
		input := strings.NewReader(fmt.Sprintf("%s\n", r))
		reader := New(bufio.NewReader(input))
		result, err = reader.Read()
		reader = nil
	}
	return result, err
}

func NewGlobalGet() SExpression {
	return &_global_get{}
}
