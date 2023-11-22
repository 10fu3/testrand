package eval

import (
	"context"
	"errors"
)

func _syntax_foreach_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}
	if arrSize != 2 {
		return CreateNil(), errors.New("malformed foreach")
	}

	//this is the list of items to iterate over
	list, err := Eval(ctx, arr[0], env)

	if err != nil {
		return CreateNil(), err
	}

	//get body
	body := arr[1]

	if !body.IsList() {
		return CreateNil(), errors.New("foreach: second argument must be a lambda")
	}

	if !body._cell._car.Equals(CreateSymbol("lambda")) {
		return CreateNil(), errors.New("foreach: second argument must be a lambda")
	}

	if !body._cell._cdr.IsList() {
		return CreateNil(), errors.New("foreach: second argument must be a lambda with a list of arguments")
	}

	if !body._cell._cdr._cell._cdr.IsList() {
		return CreateNil(), errors.New("foreach: second argument must be a lambda with a list of arguments")
	}

	bodyArg := body._cell._cdr._cell._car

	if !bodyArg.IsConsCell() {
		return CreateNil(), errors.New("foreach: second argument must be a lambda with a list of arguments")
	}

	hasParamsForIndex := !bodyArg.IsList()

	//check if list is a list
	if !list.IsList() {
		return CreateNil(), errors.New("foreach: first argument must be a list")
	}

	//get the list as an array
	listArr, listSize, err := ToArray(list)

	closure, err := Eval(ctx, body, env)

	if err != nil {
		return CreateNil(), err
	}

	if !closure.IsClosure() {
		return CreateNil(), errors.New("foreach: second argument must be a lambda")
	}

	for i := uint64(0); i < listSize; i++ {
		if hasParamsForIndex {
			_, applyFuncErr := closure._applyFunc(closure, ctx, env, CreateConsCell(listArr[i], CreateConsCell(CreateInt(int64(i)), CreateConsCell(CreateNil(), CreateNil()))))
			if applyFuncErr != nil {
				return CreateNil(), applyFuncErr
			}
		} else {
			_, applyFuncErr := closure._applyFunc(closure, ctx, env, CreateConsCell(listArr[i], CreateConsCell(CreateNil(), CreateNil())))
			if applyFuncErr != nil {
				return CreateNil(), applyFuncErr
			}
		}
	}

	return CreateNil(), nil
}

func NewForeach() *Sexpression {
	return CreateSpecialForm("foreach", _syntax_foreach_Apply)
}
