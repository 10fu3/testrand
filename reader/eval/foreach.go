package eval

import (
	"context"
	"errors"
)

type _foreach struct{}

func (_ *_foreach) TypeId() string {
	return "special_form.foreach"
}

func (_ *_foreach) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_foreach) String() string {
	return "#<syntax foreach>"
}

func (_ *_foreach) IsList() bool {
	return false
}

func (q *_foreach) Equals(sexp SExpression) bool {
	return q.TypeId() == sexp.TypeId()
}

func (_ *_foreach) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)

	if err != nil {
		return nil, err
	}
	if len(arr) != 2 {
		return nil, errors.New("malformed foreach")
	}

	//this is the list of items to iterate over
	list, err := Eval(ctx, arr[0], env)

	if err != nil {
		return nil, err
	}

	//get body
	body := arr[1]

	if !body.IsList() {
		return nil, errors.New("foreach: second argument must be a lambda")
	}

	if !body.(ConsCell).GetCar().Equals(NewSymbol("lambda")) {
		return nil, errors.New("foreach: second argument must be a lambda")
	}

	if !body.(ConsCell).GetCdr().IsList() {
		return nil, errors.New("foreach: second argument must be a lambda with a list of arguments")
	}

	if !body.(ConsCell).GetCdr().(ConsCell).GetCdr().IsList() {
		return nil, errors.New("foreach: second argument must be a lambda with a list of arguments")
	}

	bodyArg := body.(ConsCell).GetCdr().(ConsCell).GetCar()

	if bodyArg.TypeId() != "cons_cell" {
		return nil, errors.New("foreach: second argument must be a lambda with a list of arguments")
	}

	hasParamsForIndex := !bodyArg.IsList()

	//check if list is a list
	if !list.IsList() {
		return nil, errors.New("foreach: first argument must be a list")
	}

	//get the list as an array
	listArr, err := ToArray(list)

	rawClosure, err := Eval(ctx, body, env)

	if err != nil {
		return nil, err
	}

	if rawClosure.TypeId() != "closure" {
		return nil, errors.New("foreach: second argument must be a lambda")
	}

	closure := rawClosure.(*_closure)

	for i := 0; i < len(listArr); i++ {
		if hasParamsForIndex {
			_, err := closure.Apply(ctx, env, NewConsCell(listArr[i], NewConsCell(NewInt(int64(i)), NewConsCell(NewNil(), NewNil()))))
			if err != nil {
				return nil, err
			}
		} else {
			_, err := closure.Apply(ctx, env, NewConsCell(listArr[i], NewConsCell(NewNil(), NewNil())))
			if err != nil {
				return nil, err
			}
		}
	}

	return NewNil(), nil
}

func NewForeach() SExpression {
	return &_foreach{}
}
