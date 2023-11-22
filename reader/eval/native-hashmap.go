package eval

import (
	"context"
	"errors"
	"testrand/cmap"
)

//implement golang native hashmap on lisp

func _new_native_hashmap_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	return &Sexpression{
		_sexp_type_id: SexpressionTypeNativeHashmap,
		_native_map:   cmap.New[interface{}](),
	}, nil
}

func NewNativeHashmap() *Sexpression {
	return CreateSubroutine("new-native-hashmap", _new_native_hashmap_Apply)
}

func _put_native_hashmap_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	args, argsLen, err := ToArray(arguments)
	if err != nil {
		return CreateNil(), err
	}

	if 3 > argsLen {
		return CreateNil(), errors.New("need arguments size is 3")
	}

	if args[0]._sexp_type_id != SexpressionTypeNativeHashmap {
		return CreateNil(), errors.New("need arguments type is native.hashmap")
	}

	if !args[1].IsString() {
		return CreateNil(), errors.New("need arguments type is string")
	}

	args[0]._native_map.Set(args[1]._string, args[2])

	return args[2], nil
}

func NewPutNativeHashmap() *Sexpression {
	return CreateSubroutine("put-native-hashmap", _put_native_hashmap_Apply)
}

func _get_native_hashmap_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	args, argsLen, err := ToArray(arguments)
	if err != nil {
		return CreateNil(), err
	}

	if 2 > argsLen {
		return CreateNil(), errors.New("need arguments size is 2")
	}

	if !args[0].IsNativeHashmap() {
		return CreateNil(), errors.New("need arguments type is native.hashmap")
	}

	if !args[1].IsString() {
		return CreateNil(), errors.New("need arguments type is string")
	}

	v, ok := args[0]._native_map.Get(args[1]._string)

	if !ok {
		if 3 == len(args) {
			return args[2], nil
		}
		return CreateNil(), nil
	}

	convert, ok := v.(*Sexpression)

	if !ok {
		return CreateNativeValue(v), nil
	}

	return convert, nil
}

func NewGetNativeHashmap() *Sexpression {
	return CreateSubroutine("get-native-hashmap", _get_native_hashmap_Apply)
}

func _key_value_pair_foreach_native_hashmap_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	args, argsLen, err := ToArray(arguments)
	if err != nil {
		return CreateNil(), err
	}

	if 2 > argsLen {
		return CreateNil(), errors.New("need arguments size is 2")
	}

	nativeHashMap, err := Eval(ctx, args[0], env)

	if err != nil || !nativeHashMap.IsNativeHashmap() {
		return CreateNil(), errors.New("need arguments type is native.hashmap")
	}

	keyValueLambda, err := Eval(ctx, args[1], env)

	if err != nil {
		return CreateNil(), err
	}

	rawDict := nativeHashMap._native_map

	for item := range rawDict.IterBuffered() {
		key := item.Key
		value := item.Val
		convert, ok := value.(*Sexpression)
		evalTarget := CreateNil()
		if ok {
			evalTarget = CreateConsCell(keyValueLambda,
				CreateConsCell(CreateString(key), convert))
		}

		evalTarget = CreateConsCell(keyValueLambda,
			CreateConsCell(CreateString(key), CreateNativeValue(value)))

		_, evaluatedErr := Eval(ctx, evalTarget, env)
		if evaluatedErr != nil {
			return CreateNil(), evaluatedErr
		}
	}

	return CreateNil(), nil
}

func NewKeyValuePairForeachNativeHashmap() *Sexpression {
	return CreateSubroutine("key-value-pair-foreach-native-hashmap", _key_value_pair_foreach_native_hashmap_Apply)
}

func _key_value_pair_native_hashmap_to_cons_cell_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	args, argsLen, err := ToArray(arguments)
	if err != nil {
		return CreateNil(), err
	}

	if 1 > argsLen {
		return CreateNil(), errors.New("need arguments size is 1")
	}

	if !args[0].IsNativeHashmap() {
		return CreateNil(), errors.New("need arguments type is native.hashmap")
	}

	var top = CreateConsCell(CreateNil(), CreateNil())
	var consCell = top

	m := args[0]._native_map

	for item := range m.IterBuffered() {
		key := item.Key
		value := item.Val
		convert, ok := value.(*Sexpression)
		if ok {
			consCell._cell._car = CreateConsCell(CreateString(key), convert)
		} else {
			consCell._cell._car = CreateConsCell(CreateString(key), CreateNativeValue(value))
		}
		consCell._cell._cdr = CreateConsCell(CreateNil(), CreateNil())
		consCell = consCell._cell._cdr
	}

	return top, nil
}
