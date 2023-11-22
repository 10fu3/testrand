package eval

import (
	"context"
	"errors"
)

func _subr_new_native_array_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	return &Sexpression{
		_sexp_type_id: SexpressionTypeNativeArray,
		_native_arr:   make([]interface{}, 0),
	}, nil
}

func NewNativeArray() *Sexpression {
	return CreateSubroutine("new-array", _subr_new_native_array_Apply)
}

func _subr_get_native_array_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	argsArr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}

	if 2 != arrSize {
		return CreateNil(), errors.New("wrong number of arguments")
	}

	arr := argsArr[0]._native_arr

	index := argsArr[1]._number

	rawCastedArr, ok := arr.([]interface{})

	if index < 0 || index >= int64(len(rawCastedArr)) {
		return CreateNil(), errors.New("index out of range")
	}

	v, ok := rawCastedArr[index].(*Sexpression)

	if ok {
		return v, nil
	}

	return CreateNativeValue(rawCastedArr[index]), nil
}

func NewGetIndexNativeArray() *Sexpression {
	return CreateSubroutine("get-index-array", _subr_get_native_array_Apply)
}

func _set_native_array_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	argsArr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}

	if 3 != arrSize {
		return CreateNil(), errors.New("wrong number of arguments")
	}

	arr := argsArr[0]._native_arr.([]interface{})

	index := argsArr[1]._number

	if index < 0 || index >= int64(arrSize) {
		return CreateNil(), errors.New("index out of range")
	}

	arr[index] = argsArr[2]

	return argsArr[2], nil
}

func NewSetIndexNativeArray() *Sexpression {
	return CreateSubroutine("set-index-array", _set_native_array_Apply)
}

func _subr__length_native_array_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	argsArr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}

	if 1 != len(argsArr) {
		return CreateNil(), errors.New("wrong number of arguments")
	}

	_, ok := argsArr[0]._native_arr.([]interface{})

	if !ok {
		return CreateNil(), errors.New("wrong type of arguments")
	}

	return CreateInt(int64(arrSize)), nil
}

func NewLengthNativeArray() *Sexpression {
	return CreateSubroutine("length-array", _subr__length_native_array_Apply)
}

func _subr_append_native_array_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	argsArr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}

	if 2 != arrSize {
		return CreateNil(), errors.New("wrong number of arguments")
	}

	arr := argsArr[0]._native_arr.([]interface{})

	arr = append(arr, argsArr[1])

	argsArr[0]._native_arr = arr

	return argsArr[0], nil
}

func NewAppendNativeArray() *Sexpression {
	return CreateSubroutine("append-array", _subr_append_native_array_Apply)
}

func _subr_native_array_to_list_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	argsArr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}

	if 1 != arrSize {
		return CreateNil(), errors.New("wrong number of arguments")
	}

	arr, ok := argsArr[0]._native_arr.([]*Sexpression)

	if !ok {
		return CreateNil(), errors.New("wrong type of arguments")
	}

	return ToConsCell(arr), nil
}

func NewNativeArrayToList() *Sexpression {
	return CreateSubroutine("array->list", _subr_native_array_to_list_Apply)
}

func _subr_list_to_native_array_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	argsArr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}

	if 1 != arrSize {
		return CreateNil(), errors.New("wrong number of arguments")
	}

	nativeArr, _, err := ToArray(argsArr[0])

	if err != nil {
		return CreateNil(), err
	}

	return &Sexpression{
		_sexp_type_id: SexpressionTypeNativeArray,
		_native_arr:   nativeArr,
	}, nil
}

func NewListToNativeArray() *Sexpression {
	return CreateSubroutine("list->array", _subr_list_to_native_array_Apply)
}

func _syntax_foreach_native_array_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	args, argsLen, err := ToArray(arguments)
	if err != nil {
		return CreateNil(), err
	}

	if 2 > argsLen {
		return CreateNil(), errors.New("need arguments size is 2")
	}

	nativeArray, err := Eval(ctx, args[0], env)

	if err != nil || nativeArray.SexpressionTypeId() != SexpressionTypeNativeArray {
		return CreateNil(), errors.New("need arguments type is native_array")
	}

	lambda, err := Eval(ctx, args[1], env)

	if err != nil {
		return CreateNil(), err
	}

	if lambda.SexpressionTypeId() != SexpressionTypeClosure {
		return CreateNil(), errors.New("need arguments type is closure")
	}

	if lambda.GetFormalsCount() != 1 && lambda.GetFormalsCount() != 2 {
		return CreateNil(), errors.New("need arguments size is 1 or 2")
	}

	sexpConveterArr, sexpConveterr := nativeArray._native_arr.([]*Sexpression)

	var convertMode = false
	if !sexpConveterr {
		convertMode = true
	}

	if convertMode {
		nativeArr, okNativeConvert := nativeArray._native_arr.([]interface{})
		if !okNativeConvert {
			return CreateNil(), errors.New("need arguments type is native_array")
		}
		if lambda.GetFormalsCount() == 2 {
			var argsSexp = CreateEmptyList()
			for i, v := range nativeArr {
				argsSexp._cell._car = CreateInt(int64(i))
				argsSexp._cell._cdr = CreateConsCell(CreateNativeValue(v), CreateNil())
				_, err := lambda._applyFunc(lambda, ctx, env, argsSexp)

				if err != nil {
					return CreateNil(), err
				}
			}
			return CreateNil(), nil
		}
		var childArgsSexp = CreateEmptyList()
		var argsSexp = CreateConsCell(CreateNil(), childArgsSexp)
		for _, v := range nativeArr {
			argsSexp._cell._car = CreateNativeValue(v)
			_, lambdaRunErr := lambda._applyFunc(lambda, ctx, env, argsSexp)

			if lambdaRunErr != nil {
				return CreateNil(), lambdaRunErr
			}
		}
	}
	if lambda.GetFormalsCount() == 2 {
		var argsSexp = CreateEmptyList()
		for i, v := range sexpConveterArr {
			argsSexp._cell._car = CreateInt(int64(i))
			argsSexp._cell._cdr = CreateConsCell(v, CreateNil())
			_, err := lambda._applyFunc(lambda, ctx, env, argsSexp)

			if err != nil {
				return CreateNil(), err
			}
		}
		return CreateNil(), nil
	}
	var childArgsSexp = CreateEmptyList()
	var argsSexp = CreateConsCell(CreateNil(), childArgsSexp)
	for _, v := range sexpConveterArr {
		argsSexp._cell._car = v
		_, lambdaRunErr := lambda._applyFunc(lambda, ctx, env, argsSexp)

		if lambdaRunErr != nil {
			return CreateNil(), lambdaRunErr
		}
	}
	return CreateNil(), nil
}

func NewForeachNativeArray() *Sexpression {
	return CreateSpecialForm("foreach-array", _syntax_foreach_native_array_Apply)
}
