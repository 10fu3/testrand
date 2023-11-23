package eval

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"sync"
)

//implement golang native hashmap on lisp

type _native_hashmap struct {
	id string
	M  map[string]SExpression
	sync.Mutex
}

func (_ *_native_hashmap) TypeId() string {
	return "native.hashmap"
}

func (_ *_native_hashmap) AtomId() SExpressionType {
	return SExpressionTypeNativeHashmap
}

func (_ *_native_hashmap) String() string {
	return "#<native hashmap>"
}

func (_ *_native_hashmap) IsList() bool {
	return false
}

func (l *_native_hashmap) Equals(sexp SExpression) bool {
	if sexp.TypeId() != "native.hashmap" {
		return false
	}

	return l.id == sexp.(*_native_hashmap).id
}

type _new_native_hashmap struct{}

func (_ *_new_native_hashmap) TypeId() string {
	return "subroutine.new-native-hashmap"
}

func (_ *_new_native_hashmap) AtomId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_new_native_hashmap) String() string {
	return "#<subr new-native-hashmap>"
}

func (_ *_new_native_hashmap) IsList() bool {
	return false
}

func (l *_new_native_hashmap) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_new_native_hashmap) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	return &_native_hashmap{
		id: uuid.NewString(),
		M:  make(map[string]SExpression),
	}, nil
}

func NewNativeHashmap() SExpression {
	return &_new_native_hashmap{}
}

type _put_native_hashmap struct{}

func (_ *_put_native_hashmap) TypeId() string {
	return "subroutine.put-native-hashmap"
}

func (_ *_put_native_hashmap) AtomId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_put_native_hashmap) String() string {
	return "#<subr put-native-hashmap>"
}

func (_ *_put_native_hashmap) IsList() bool {
	return false
}

func (l *_put_native_hashmap) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_put_native_hashmap) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 3 > argsLength {
		return nil, errors.New("need arguments size is 3")
	}

	if args[0].TypeId() != "native.hashmap" {
		return nil, errors.New("need arguments type is native.hashmap")
	}

	if args[1].TypeId() != "string" {
		return nil, errors.New("need arguments type is string")
	}

	args[0].(*_native_hashmap).Lock()
	args[0].(*_native_hashmap).M[args[1].(Str).GetValue()] = args[2]
	args[0].(*_native_hashmap).Unlock()

	return args[2], nil
}

func NewPutNativeHashmap() SExpression {
	return &_put_native_hashmap{}
}

type _get_native_hashmap struct{}

func (_ *_get_native_hashmap) TypeId() string {
	return "subroutine.get-native-hashmap"
}

func (_ *_get_native_hashmap) AtomId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_get_native_hashmap) String() string {
	return "#<subr get-native-hashmap>"
}

func (_ *_get_native_hashmap) IsList() bool {
	return false
}

func (l *_get_native_hashmap) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_get_native_hashmap) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 2 > argsLength {
		return nil, errors.New("need arguments size is 2")
	}

	if args[0].TypeId() != "native.hashmap" {
		return nil, errors.New("need arguments type is native.hashmap")
	}

	if args[1].TypeId() != "string" {
		return nil, errors.New("need arguments type is string")
	}

	args[0].(*_native_hashmap).Lock()
	v, ok := args[0].(*_native_hashmap).M[args[1].(Str).GetValue()]
	args[0].(*_native_hashmap).Unlock()

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

type _key_value_pair_foreach_native_hashmap struct{}

func (_ *_key_value_pair_foreach_native_hashmap) TypeId() string {
	return "special_form.kv-set-native-hashmap"
}

func (_ *_key_value_pair_foreach_native_hashmap) AtomId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_key_value_pair_foreach_native_hashmap) String() string {
	return "#<syntax kv-set-native-hashmap>"
}

func (_ *_key_value_pair_foreach_native_hashmap) IsList() bool {
	return false
}

func (l *_key_value_pair_foreach_native_hashmap) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_key_value_pair_foreach_native_hashmap) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 2 > argsLength {
		return nil, errors.New("need arguments size is 2")
	}

	nativeHashMap, err := Eval(ctx, args[0], env)

	if err != nil || nativeHashMap.TypeId() != "native.hashmap" {
		return nil, errors.New("need arguments type is native.hashmap")
	}

	keyValueLambda, err := Eval(ctx, args[1], env)

	if err != nil {
		return nil, err
	}

	nativeHashMap.(*_native_hashmap).Lock()
	//for loop key value
	for key, value := range nativeHashMap.(*_native_hashmap).M {
		evalTarget := NewConsCell(keyValueLambda,
			NewConsCell(NewString(key), value))

		Eval(ctx, evalTarget, env)
	}
	nativeHashMap.(*_native_hashmap).Unlock()

	return NewNil(), nil
}

func NewKeyValuePairNativeHashmap() SExpression {
	return &_key_value_pair_foreach_native_hashmap{}
}

type _key_value_pair_native_hashmap_to_cons_cell struct{}

func (_ *_key_value_pair_native_hashmap_to_cons_cell) TypeId() string {
	return "subroutine.kv-native-hashmap"
}

func (_ *_key_value_pair_native_hashmap_to_cons_cell) AtomId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_key_value_pair_native_hashmap_to_cons_cell) String() string {
	return "#<subr kv-native-hashmap>"
}

func (_ *_key_value_pair_native_hashmap_to_cons_cell) IsList() bool {
	return false
}

func (l *_key_value_pair_native_hashmap_to_cons_cell) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_key_value_pair_native_hashmap_to_cons_cell) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 1 > argsLength {
		return nil, errors.New("need arguments size is 1")
	}

	if args[0].TypeId() != "native.hashmap" {
		return nil, errors.New("need arguments type is native.hashmap")
	}

	args[0].(*_native_hashmap).Lock()
	var top = NewConsCell(NewNil(), NewNil()).(*_cons_cell)
	var consCell SExpression = top
	for key, value := range args[0].(*_native_hashmap).M {
		consCell.(*_cons_cell).Car = (NewConsCell(NewString(key), value))
		consCell.(*_cons_cell).Cdr = (NewConsCell(NewNil(), NewNil()))
		consCell = consCell.(*_cons_cell).GetCdr()
	}
	args[0].(*_native_hashmap).Unlock()

	return top, nil
}

func NewKeyValuePairNativeHashmapToConsCell() SExpression {
	return &_key_value_pair_native_hashmap_to_cons_cell{}
}
