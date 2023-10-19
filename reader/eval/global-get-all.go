package eval

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strings"
)

type _global_get_all struct {
}

func (_ *_global_get_all) Type() string {
	return "special_form.global_get_all"
}

func (_ *_global_get_all) String() string {
	return "#<syntax global_get_all>"
}

func (_ *_global_get_all) IsList() bool {
	return false
}

func (l *_global_get_all) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_global_get_all) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	if "cons_cell" != args.Type() {
		return nil, errors.New("type error")
	}

	var err error

	if err != nil {
		return nil, err
	}

	res, err := env.GetSuperGlobalEnv().GetAll()

	if err != nil {
		return nil, err
	}

	var list = &_cons_cell{Car: NewNil(), Cdr: NewNil()}

	for _, kv := range res {
		input := strings.NewReader(fmt.Sprintf("%s\n", kv.Value))
		reader := New(bufio.NewReader(input))
		result, err := reader.Read()
		if err != nil {
			continue
		}
		keyValue := NewConsCell(NewSymbol(kv.Key), result)
		list.Car = keyValue
		list.Cdr = NewConsCell(NewNil(), NewNil())
		list = list.Cdr.(*_cons_cell)
	}
	return list, nil
}

func NewGlobalGetAll() SExpression {
	return &_global_get_all{}
}
