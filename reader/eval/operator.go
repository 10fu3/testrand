package eval

import (
	"context"
	"errors"
)

type _is_not struct {
}

func (_ *_is_not) TypeId() string {
	return "subroutine.is_equals"
}

func (_ *_is_not) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_is_not) String() string {
	return "#<subr eq?>"
}

func (_ *_is_not) IsList() bool {
	return false
}

func (i *_is_not) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_is_not) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 1 {
		return nil, errors.New("malformed not")
	}

	first := args[0]

	if "bool" != first.TypeId() {
		return first, nil
	}

	return NewBool(!first.(Bool).GetValue()), nil
}

func NewIsNot() SExpression {
	return &_is_not{}
}

type _or struct {
}

func (_ _or) TypeId() string {
	return "special_form.or"
}

func (_ _or) AtomId() AtomType {
	return AtomTypeSpecialForm
}

func (_ _or) String() string {
	return "#<syntax #or>"
}

func (_ _or) IsList() bool {
	return false
}

func (a _or) Equals(sexp SExpression) bool {
	return a.TypeId() == sexp.TypeId()
}

func (_ _or) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	evaluatedElm := NewConsCell(NewNil(), NewNil()).(SExpression)
	var err error

	for i := uint64(0); i < argsLength; i++ {
		evaluatedElm, err = Eval(ctx, args[i], env)
		if err != nil {
			return nil, err
		}
		if !NewBool(false).Equals(evaluatedElm) {
			return evaluatedElm, nil
		}
	}

	return evaluatedElm, nil
}

func NewOr() SExpression {
	return &_or{}
}

type _is_equals struct {
}

func (_ *_is_equals) TypeId() string {
	return "subroutine.is_equals"
}

func (_ *_is_equals) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_is_equals) String() string {
	return "#<subr eq?>"
}

func (_ *_is_equals) IsList() bool {
	return false
}

func (i *_is_equals) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_is_equals) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("malformed if syntax")
	}

	first := args[0]

	second := args[1]

	return NewBool(first.Equals(second)), nil
}

func NewIsEq() SExpression {
	return &_is_equals{}
}

type _is_greater_than struct{}

func (_ *_is_greater_than) TypeId() string {
	return "subroutine.is_greater_than"
}

func (_ *_is_greater_than) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_is_greater_than) String() string {
	return "#<subr > >"
}

func (_ *_is_greater_than) IsList() bool {
	return false
}

func (i *_is_greater_than) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_is_greater_than) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("malformed if syntax")
	}

	first := args[0]

	second := args[1]

	return NewBool(first.(Number).GetValue() > second.(Number).GetValue()), nil
}

func NewIsGreaterThan() SExpression {
	return &_is_greater_than{}
}

type _is_less_than struct{}

func (_ *_is_less_than) TypeId() string {
	return "subroutine.is_less_than"
}

func (_ *_is_less_than) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_is_less_than) String() string {
	return "#<subr < >"
}

func (_ *_is_less_than) IsList() bool {
	return false
}

func (i *_is_less_than) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_is_less_than) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("malformed if syntax")
	}

	first := args[0]

	second := args[1]

	return NewBool(first.(Number).GetValue() < second.(Number).GetValue()), nil
}

func NewIsLessThan() SExpression {
	return &_is_less_than{}
}

type _is_greater_than_or_equal struct{}

func (_ *_is_greater_than_or_equal) TypeId() string {
	return "subroutine.is_greater_than_or_equal"
}

func (_ *_is_greater_than_or_equal) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_is_greater_than_or_equal) String() string {
	return "#<subr >= >"
}

func (_ *_is_greater_than_or_equal) IsList() bool {
	return false
}

func (i *_is_greater_than_or_equal) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_is_greater_than_or_equal) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("malformed if syntax")
	}

	first := args[0]

	second := args[1]

	return NewBool(first.(Number).GetValue() >= second.(Number).GetValue()), nil
}

func NewIsGreaterThanOrEqual() SExpression {
	return &_is_greater_than_or_equal{}
}

type _is_less_than_or_equal struct{}

func (_ *_is_less_than_or_equal) TypeId() string {
	return "subroutine.is_less_than_or_equal"
}

func (_ *_is_less_than_or_equal) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_is_less_than_or_equal) String() string {
	return "#<subr <= >"
}

func (_ *_is_less_than_or_equal) IsList() bool {
	return false
}

func (i *_is_less_than_or_equal) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_is_less_than_or_equal) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("malformed if syntax")
	}

	first := args[0]

	second := args[1]

	return NewBool(first.(Number).GetValue() <= second.(Number).GetValue()), nil
}

func NewIsLessThanOrEqual() SExpression {
	return &_is_less_than_or_equal{}
}

type _is_null struct{}

func (_ *_is_null) TypeId() string {
	return "subroutine.is_null"
}

func (_ *_is_null) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_is_null) String() string {
	return "#<subr null?>"
}

func (_ *_is_null) IsList() bool {
	return false
}

func (i *_is_null) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_is_null) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 1 {
		return nil, errors.New("malformed if syntax")
	}

	first := args[0]

	return NewBool(AtomTypeNil == first.AtomId()), nil
}

func NewIsNull() SExpression {
	return &_is_null{}
}

type _is_num_equal struct{}

func (_ *_is_num_equal) TypeId() string {
	return "subroutine.is_equal"
}

func (_ *_is_num_equal) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_is_num_equal) String() string {
	return "#<subr =>"
}

func (_ *_is_num_equal) IsList() bool {
	return false
}

func (i *_is_num_equal) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_is_num_equal) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("malformed if syntax")
	}

	first := args[0]

	second := args[1]

	return NewBool(first.(Number).GetValue() == second.(Number).GetValue()), nil
}

func NewIsNumEqual() SExpression {
	return &_is_num_equal{}
}
