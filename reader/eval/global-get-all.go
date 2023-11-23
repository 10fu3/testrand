package eval

import (
	"bufio"
	"context"
	"fmt"
	"strings"
)

type _global_get_all struct {
}

func (_ *_global_get_all) TypeId() string {
	return "special_form.global_get_all"
}

func (_ *_global_get_all) AtomId() AtomType {
	return AtomTypeSpecialForm
}

func (_ *_global_get_all) String() string {
	return "#<syntax global_get_all>"
}

func (_ *_global_get_all) IsList() bool {
	return false
}

func (l *_global_get_all) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_global_get_all) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	res, err := env.GetSuperGlobalEnv().GetAll()

	if err != nil {
		return nil, err
	}

	var list = &_cons_cell{Car: NewNil(), Cdr: NewNil()}
	var parent = list
	for _, kv := range res {
		input := strings.NewReader(fmt.Sprintf("%s\n", kv.Value))
		reader := New(bufio.NewReader(input))
		result, err := reader.Read()
		if err != nil {
			continue
		}
		_, keyName, foundName := strings.Cut(kv.Key, fmt.Sprintf("/env/%s/", env.GetParentId()))
		if !foundName {
			continue
		}
		keyValue := NewConsCell(NewSymbol(keyName), result)
		list.Car = keyValue
		list.Cdr = NewConsCell(NewNil(), NewNil())
		list = list.Cdr.(*_cons_cell)
	}
	return parent, nil
}

func NewGlobalGetAll() SExpression {
	return &_global_get_all{}
}
