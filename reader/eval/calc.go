package eval

import (
	"context"
	"errors"
)

type _add struct{}

func (_ *_add) TypeId() string {
	return "subroutine.add"
}

func (_ *_add) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_add) String() string {
	return "#<subr +>"
}

func (_ *_add) IsList() bool {
	return false
}

func (a *_add) Equals(sexp SExpression) bool {
	return a.TypeId() == sexp.TypeId()
}

func (_ *_add) Apply(ctx context.Context, _ Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 0 == argsLength {
		return NewInt(0), nil
	}

	arrIndex := uint64(1)

	var result Number = nil
	if "number" != args[0].TypeId() {
		return nil, errors.New("need arguments type is number, but got " + args[0].TypeId())
	}
	result = args[0].(Number)
	for ; arrIndex < argsLength; arrIndex++ {
		if "number" != args[arrIndex].TypeId() {
			return nil, errors.New("need arguments type is number, but got " + args[arrIndex].TypeId())
		}
		result = NewInt(result.GetValue() + args[arrIndex].(Number).GetValue())
	}
	return result, nil
}

func NewAdd() SExpression {
	return &_add{}
}

type _minus struct{}

func (_ *_minus) TypeId() string {
	return "subroutine.minus"
}

func (_ *_minus) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_minus) String() string {
	return "#<subr ->"
}

func (_ *_minus) IsList() bool {
	return false
}

func (a *_minus) Equals(sexp SExpression) bool {
	return a.TypeId() == sexp.TypeId()
}

func (_ *_minus) Apply(ctx context.Context, _ Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("need arguments size is 2")
	}

	arrIndex := uint64(1)

	var result Number = nil
	if "number" != args[0].TypeId() {
		return nil, errors.New("need arguments type is number, but got " + args[0].TypeId())
	}
	result = args[0].(Number)
	for ; arrIndex < argsLength; arrIndex++ {
		if "number" != args[arrIndex].TypeId() {
			return nil, errors.New("need arguments type is number, but got " + args[arrIndex].TypeId())
		}
		result = NewInt(result.GetValue() - args[arrIndex].(Number).GetValue())
	}
	return result, nil
}

func NewMinus() SExpression {
	return &_minus{}
}

type _multiply struct{}

func (_ *_multiply) TypeId() string {
	return "subroutine.multiply"
}

func (_ *_multiply) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_multiply) String() string {
	return "#<subr *>"
}

func (_ *_multiply) IsList() bool {
	return false
}

func (a *_multiply) Equals(sexp SExpression) bool {
	return a.TypeId() == sexp.TypeId()
}

func (_ *_multiply) Apply(ctx context.Context, _ Environment, arr []SExpression, argsLength uint64) (SExpression, error) {

	if 0 == argsLength {
		return NewInt(0), nil
	}

	arrIndex := uint64(1)

	var result Number = nil
	if "number" != arr[0].TypeId() {
		return nil, errors.New("need arguments type is number, but got " + arr[0].TypeId())
	}
	result = arr[0].(Number)
	for ; arrIndex < argsLength; arrIndex++ {
		if "number" != arr[arrIndex].TypeId() {
			return nil, errors.New("need arguments type is number, but got " + arr[arrIndex].TypeId())
		}
		result = NewInt(result.GetValue() * arr[arrIndex].(Number).GetValue())
	}
	return result, nil
}

func NewMultiply() SExpression {
	return &_multiply{}
}

type _divide struct{}

func (_ *_divide) TypeId() string {
	return "subroutine.divide"
}

func (_ *_divide) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_divide) String() string {
	return "#<subr />"
}

func (_ *_divide) IsList() bool {
	return false
}

func (a *_divide) Equals(sexp SExpression) bool {
	return a.TypeId() == sexp.TypeId()
}

func (_ *_divide) Apply(ctx context.Context, _ Environment, arr []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("need arguments size is 2")
	}

	arrIndex := uint64(1)

	var result Number = nil
	if "number" != arr[0].TypeId() {
		return nil, errors.New("need arguments type is number, but got " + arr[0].TypeId())
	}
	result = arr[0].(Number)
	for ; arrIndex < argsLength; arrIndex++ {
		if "number" != arr[arrIndex].TypeId() {
			return nil, errors.New("need arguments type is number, but got " + arr[arrIndex].TypeId())
		}
		result = NewInt(result.GetValue() / arr[arrIndex].(Number).GetValue())
	}
	return result, nil
}

func NewDivide() SExpression {
	return &_divide{}
}

type _mod struct{}

func (_ *_mod) TypeId() string {
	return "subroutine.mod"
}

func (_ *_mod) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_mod) String() string {
	return "#<subr %>"
}

func (_ *_mod) IsList() bool {
	return false
}

func (a *_mod) Equals(sexp SExpression) bool {
	return a.TypeId() == sexp.TypeId()
}

func (_ *_mod) Apply(ctx context.Context, _ Environment, arr []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("need arguments size is 2")
	}

	if "number" != arr[0].TypeId() {
		return nil, errors.New("need arguments type is number, but got " + arr[0].TypeId())
	}

	if "number" != arr[1].TypeId() {
		return nil, errors.New("need arguments type is number, but got " + arr[1].TypeId())
	}

	base := arr[0].(Number).GetValue()
	target := arr[1].(Number).GetValue()

	return NewInt(base % target), nil
}

func NewMod() SExpression {
	return &_mod{}
}
