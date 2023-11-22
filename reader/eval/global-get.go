package eval

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"
	"strings"
)

func _syntax_global_get_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	if !args.IsConsCell() {
		return CreateNil(), errors.New("type error")
	}

	cell := args._cell

	name := cell._car._symbol

	defaultArg := (func() *Sexpression {
		if IsEmptyList(cell._cdr) {
			return CreateNil()
		}
		return cell._cdr._cell._car
	})()

	var err error

	if err != nil {
		return CreateNil(), err
	}

	var result *Sexpression
	if ctx.Value("transaction") != nil {
		transaction := ctx.Value("transaction").(concurrency.STM)
		existKey := transaction.Rev(fmt.Sprintf("/env/%s/%s", env._env_parentId, name.String()))

		if !defaultArg.IsNil() && existKey == 0 {
			return defaultArg, nil
		}

		var r = transaction.Get(fmt.Sprintf("/env/%s/%s", env._env_parentId, name.String()))
		input := strings.NewReader(fmt.Sprintf("%s\n", r))
		rd := New(bufio.NewReader(input))
		result, err = rd.Read()
		rd = nil

	} else {
		r, globalGetErr := env._env_superGlobalEnv.Get(fmt.Sprintf("/env/%s/%s", env._env_parentId, name.String()))
		if globalGetErr != nil {
			return CreateNil(), globalGetErr
		}
		input := strings.NewReader(fmt.Sprintf("%s\n", r))
		rd := New(bufio.NewReader(input))
		result, err = rd.Read()
		rd = nil
	}
	return result, err
}

func NewGlobalGet() *Sexpression {
	return CreateSpecialForm("global-get", _syntax_global_get_Apply)
}
