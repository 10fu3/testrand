package eval

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

//implement golang native hashmap on lisp

type _native_hashmap struct {
	id string
	M  map[string]SExpression
}

func (_ *_native_hashmap) Type() string {
	return "native.hashmap"
}

func (_ *_native_hashmap) String() string {
	return "#<native hashmap>"
}

func (_ *_native_hashmap) IsList() bool {
	return false
}

func (l *_native_hashmap) Equals(sexp SExpression) bool {
	if sexp.Type() != "native.hashmap" {
		return false
	}

	return l.id == sexp.(*_native_hashmap).id
}

type _new_native_hashmap struct{}

func (_ *_new_native_hashmap) Type() string {
	return "subroutine.new-native-hashmap"
}

func (_ *_new_native_hashmap) String() string {
	return "#<subr new-native-hashmap>"
}

func (_ *_new_native_hashmap) IsList() bool {
	return false
}

func (l *_new_native_hashmap) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_new_native_hashmap) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {
	return &_native_hashmap{
		id: uuid.NewString(),
		M:  make(map[string]SExpression),
	}, nil
}

func NewNativeHashmap() SExpression {
	return &_new_native_hashmap{}
}

type _put_native_hashmap struct{}

func (_ *_put_native_hashmap) Type() string {
	return "subroutine.put-native-hashmap"
}

func (_ *_put_native_hashmap) String() string {
	return "#<subr put-native-hashmap>"
}

func (_ *_put_native_hashmap) IsList() bool {
	return false
}

func (l *_put_native_hashmap) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_put_native_hashmap) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {
	args, err := ToArray(arguments)
	if err != nil {
		return nil, err
	}

	if 3 > len(args) {
		return nil, errors.New("need arguments size is 3")
	}

	if args[0].Type() != "native.hashmap" {
		return nil, errors.New("need arguments type is native.hashmap")
	}

	if args[1].Type() != "string" {
		return nil, errors.New("need arguments type is string")
	}

	args[0].(*_native_hashmap).M[args[1].(Str).GetValue()] = args[2]

	return args[2], nil
}

func NewPutNativeHashmap() SExpression {
	return &_put_native_hashmap{}
}

type _get_native_hashmap struct{}

func (_ *_get_native_hashmap) Type() string {
	return "subroutine.get-native-hashmap"
}

func (_ *_get_native_hashmap) String() string {
	return "#<subr get-native-hashmap>"
}

func (_ *_get_native_hashmap) IsList() bool {
	return false
}

func (l *_get_native_hashmap) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_get_native_hashmap) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {
	args, err := ToArray(arguments)
	if err != nil {
		return nil, err
	}

	if 2 > len(args) {
		return nil, errors.New("need arguments size is 2")
	}

	if args[0].Type() != "native.hashmap" {
		return nil, errors.New("need arguments type is native.hashmap")
	}

	if args[1].Type() != "string" {
		return nil, errors.New("need arguments type is string")
	}

	v, ok := args[0].(*_native_hashmap).M[args[1].(Str).GetValue()]

	if !ok {
		if 3 == len(args) {
			return args[2], nil
		}
		return NewNil(), nil
	}

	return v, nil
}

func NewGetNativeHashmap() SExpression {
	return &_get_native_hashmap{}
}

type _key_value_pair_native_hashmap struct{}

func (_ *_key_value_pair_native_hashmap) Type() string {
	return "subroutine.kv-set-native-hashmap"
}

func (_ *_key_value_pair_native_hashmap) String() string {
	return "#<subr kv-set-native-hashmap>"
}

func (_ *_key_value_pair_native_hashmap) IsList() bool {
	return false
}

func (l *_key_value_pair_native_hashmap) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_key_value_pair_native_hashmap) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {
	args, err := ToArray(arguments)
	if err != nil {
		return nil, err
	}

	if 2 > len(args) {
		return nil, errors.New("need arguments size is 2")
	}

	if args[0].Type() != "native.hashmap" {
		return nil, errors.New("need arguments type is native.hashmap")
	}

	keyValueLambda := args[1]

	//for loop key value
	for key, value := range args[0].(*_native_hashmap).M {
		evalTarget := NewConsCell(keyValueLambda,
			NewConsCell(NewString(key), value))

		Eval(ctx, evalTarget, env)
	}

	return NewNil(), nil
}

func NewKeyValuePairNativeHashmap() SExpression {
	return &_key_value_pair_native_hashmap{}
}
