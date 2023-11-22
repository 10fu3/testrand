package eval

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strings"
)

func _subr_global_get_all_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	if !args.IsConsCell() {
		return CreateNil(), errors.New("type error")
	}

	var err error

	if err != nil {
		return CreateNil(), err
	}

	res, err := env._env_superGlobalEnv.GetAll()

	if err != nil {
		return CreateNil(), err
	}

	var parent = CreateEmptyList()
	var list = parent
	for _, kv := range res {
		cell := list._cell
		input := strings.NewReader(fmt.Sprintf("%s\n", kv.Value))
		reader := New(bufio.NewReader(input))
		result, err := reader.Read()
		if err != nil {
			continue
		}
		_, keyName, foundName := strings.Cut(kv.Key, fmt.Sprintf("/env/%s/", env._env_parentId))
		if !foundName {
			continue
		}
		keyValue := CreateConsCell(GetSymbol(keyName), result)
		cell._car = keyValue
		cell._cdr = CreateConsCell(CreateNil(), CreateNil())
		list = cell._cdr
	}
	return parent, nil
}

func NewGlobalGetAll() *Sexpression {
	return CreateSubroutine("global-get-all", _subr_global_get_all_Apply)
}
