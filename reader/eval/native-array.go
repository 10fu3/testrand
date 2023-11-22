package eval

import (
	"context"
	"errors"
	"fmt"
)

type _native_array struct {
	Arr []SExpression
}

func (_ *_native_array) TypeId() string {
	return "native_array"
}

func (_ *_native_array) SExpressionTypeId() SExpressionType {
	return SExpressionTypeNativeArray
}

func (l *_native_array) String() string {
	//printout element
	var str string
	var arrayLength = len(l.Arr)
	for i, v := range l.Arr {
		str += v.String()
		if i != arrayLength-1 {
			str += ","
		}
	}
	return fmt.Sprintf("[%s]", str)
}

func (_ *_native_array) IsList() bool {
	return false
}

func (_ *_native_array) Equals(sexp SExpression) bool {
	return false
}

type _new_native_array struct{}

func (_ *_new_native_array) TypeId() string {
	return "subroutine.new-native-array"
}

func (_ *_new_native_array) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_new_native_array) String() string {
	return "#<subr new-native-array>"
}

func (_ *_new_native_array) IsList() bool {
	return false
}

func (_ *_new_native_array) Equals(sexp SExpression) bool {
	return sexp.TypeId() == "subroutine.new-native-array"
}

func (_ *_new_native_array) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	return &_native_array{
		Arr: make([]SExpression, 0),
	}, nil
}

func NewNativeArray() SExpression {
	return &_new_native_array{}
}

type _get_native_array struct{}

func (_ *_get_native_array) TypeId() string {
	return "subroutine.get-native-array"
}

func (_ *_get_native_array) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_get_native_array) String() string {
	return "#<subr get-native-array>"
}

func (_ *_get_native_array) IsList() bool {
	return false
}

func (_ *_get_native_array) Equals(sexp SExpression) bool {
	return sexp.TypeId() == "subroutine.get-native-array"
}

func (_ *_get_native_array) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 2 != argsLength {
		return nil, errors.New("wrong number of arguments")
	}

	arr := args[0].(*_native_array)

	index := args[1].(Number).GetValue()

	if index < 0 || index >= int64(len(arr.Arr)) {
		return nil, errors.New("index out of range")
	}

	return arr.Arr[index], nil
}

func NewGetNativeArray() SExpression {
	return &_get_native_array{}
}

type _set_native_array struct{}

func (_ *_set_native_array) TypeId() string {
	return "subroutine.set-native-array"
}

func (_ *_set_native_array) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_set_native_array) String() string {
	return "#<subr set-native-array>"
}

func (_ *_set_native_array) IsList() bool {
	return false
}

func (_ *_set_native_array) Equals(sexp SExpression) bool {
	return sexp.TypeId() == "subroutine.set-native-array"
}

func (_ *_set_native_array) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 3 != len(args) {
		return nil, errors.New("wrong number of arguments")
	}

	arr := args[0].(*_native_array)

	index := args[1].(Number).GetValue()

	if index < 0 || index >= int64(len(arr.Arr)) {
		return nil, errors.New("index out of range")
	}

	arr.Arr[index] = args[2]

	return args[2], nil
}

func NewSetNativeArray() SExpression {
	return &_set_native_array{}
}

type _length_native_array struct{}

func (_ *_length_native_array) TypeId() string {
	return "subroutine.length-native-array"
}

func (_ *_length_native_array) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_length_native_array) String() string {
	return "#<subr length-native-array>"
}

func (_ *_length_native_array) IsList() bool {
	return false
}

func (_ *_length_native_array) Equals(sexp SExpression) bool {
	return sexp.TypeId() == "subroutine.length-native-array"
}

func (_ *_length_native_array) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 1 != argsLength {
		return nil, errors.New("wrong number of arguments")
	}

	arr := args[0].(*_native_array)

	return NewInt(int64(len(arr.Arr))), nil
}

func NewLengthNativeArray() SExpression {
	return &_length_native_array{}
}

type _append_native_array struct{}

func (_ *_append_native_array) TypeId() string {
	return "subroutine.append-native-array"
}

func (_ *_append_native_array) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_append_native_array) String() string {
	return "#<subr append-native-array>"
}

func (_ *_append_native_array) IsList() bool {
	return false
}

func (_ *_append_native_array) Equals(sexp SExpression) bool {
	return sexp.TypeId() == "subroutine.append-native-array"
}

func (_ *_append_native_array) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 2 != argsLength {
		return nil, errors.New("wrong number of arguments")
	}

	arr := args[0].(*_native_array)

	arr.Arr = append(arr.Arr, args[1])

	return arr, nil
}

func NewAppendNativeArray() SExpression {
	return &_append_native_array{}
}

type _native_array_to_list struct{}

func (_ *_native_array_to_list) TypeId() string {
	return "subroutine.native-array-to-list"
}

func (_ *_native_array_to_list) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_native_array_to_list) String() string {
	return "#<subr native-array-to-list>"
}

func (_ *_native_array_to_list) IsList() bool {
	return false
}

func (_ *_native_array_to_list) Equals(sexp SExpression) bool {
	return sexp.TypeId() == "subroutine.native-array-to-list"
}

func (_ *_native_array_to_list) Apply(ctx context.Context, env Environment, argsArr []SExpression, argsLength uint64) (SExpression, error) {

	if 1 != argsLength {
		return nil, errors.New("wrong number of arguments")
	}

	arr := argsArr[0].(*_native_array)

	return ToConsCell(arr.Arr, uint64(len(arr.Arr))), nil
}

func NewNativeArrayToList() SExpression {
	return &_native_array_to_list{}
}

type _list_to_native_array struct{}

func (_ *_list_to_native_array) TypeId() string {
	return "subroutine.list-to-native-array"
}

func (_ *_list_to_native_array) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_list_to_native_array) String() string {
	return "#<subr list-to-native-array>"
}

func (_ *_list_to_native_array) IsList() bool {
	return false
}

func (_ *_list_to_native_array) Equals(sexp SExpression) bool {
	return sexp.TypeId() == "subroutine.list-to-native-array"
}

func (_ *_list_to_native_array) Apply(ctx context.Context, env Environment, argsArr []SExpression, argsLength uint64) (SExpression, error) {

	if 1 != argsLength {
		return nil, errors.New("wrong number of arguments")
	}

	consCell, _, err := ToArray(argsArr[0])

	if err != nil {
		return nil, err
	}

	return &_native_array{
		Arr: consCell,
	}, nil
}

func NewListToNativeArray() SExpression {
	return &_list_to_native_array{}
}

type _foreach_native_array struct{}

func (_ *_foreach_native_array) TypeId() string {
	return "special_form.foreach-native-array"
}

func (_ *_foreach_native_array) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_foreach_native_array) String() string {
	return "#<syntax foreach-native-array>"
}

func (_ *_foreach_native_array) IsList() bool {
	return false
}

func (l *_foreach_native_array) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_foreach_native_array) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 2 > argsLength {
		return nil, errors.New("need arguments size is 2")
	}

	nativeArray, err := Eval(ctx, args[0], env)

	if err != nil || nativeArray.SExpressionTypeId() != SExpressionTypeNativeArray {
		return nil, errors.New("need arguments type is native_array")
	}

	rawLambda, err := Eval(ctx, args[1], env)

	if err != nil {
		return nil, err
	}

	if rawLambda.SExpressionTypeId() != SExpressionTypeClosure {
		return nil, errors.New("need arguments type is closure")
	}

	lambda := rawLambda.(Closure)

	if lambda.GetFormalsCount() != 1 && lambda.GetFormalsCount() != 2 {
		return nil, errors.New("need arguments size is 1 or 2")
	}

	if lambda.GetFormalsCount() == 2 {
		for i, v := range nativeArray.(*_native_array).Arr {
			_, err := lambda.Apply(ctx, env, []SExpression{NewInt(int64(i)), v}, 2)

			if err != nil {
				return nil, err
			}
		}
		return NewNil(), nil
	} else {
		for _, v := range nativeArray.(*_native_array).Arr {
			_, err := lambda.Apply(ctx, env, []SExpression{v}, 1)

			if err != nil {
				return nil, err
			}
		}
		return NewNil(), nil
	}
}

func NewForeachNativeArray() SExpression {
	return &_foreach_native_array{}
}
