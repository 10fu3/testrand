package eval

import (
	"context"
	"errors"
)

func _sub_add_Apply(self *Sexpression, ctx context.Context, _ *Sexpression, args *Sexpression) (*Sexpression, error) {
	if args._sexp_type_id != SexpressionTypeConsCell {
		return CreateNil(), errors.New("need arguments")
	}

	if IsEmptyList(args) {
		return CreateInt(0), nil
	}

	arr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}
	arrIndex := uint64(1)

	var result = int64(0)
	if arr[0]._sexp_type_id != SexpressionTypeNumber {
		return CreateNil(), errors.New("need arguments type is number, but got " + arr[0]._sexp_type_id.String())
	}
	result = arr[0]._number
	for ; arrIndex < arrSize; arrIndex++ {
		if !arr[arrIndex].IsNumber() {
			return CreateNil(), errors.New("need arguments type is number, but got " + arr[arrIndex]._sexp_type_id.String())
		}
		result = result + arr[arrIndex]._number
	}
	return CreateInt(result), nil
}

func NewAdd() *Sexpression {
	return CreateSubroutine("add", _sub_add_Apply)
}

func _sub_minus_Apply(self *Sexpression, ctx context.Context, _ *Sexpression, args *Sexpression) (*Sexpression, error) {
	if args._sexp_type_id != SexpressionTypeConsCell {
		return CreateNil(), errors.New("need arguments")
	}

	if IsEmptyList(args) {
		return CreateInt(0), nil
	}

	arr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}
	arrIndex := uint64(1)

	var result = int64(0)
	if !arr[0].IsNumber() {
		return CreateNil(), errors.New("need arguments type is number, but got " + arr[0]._sexp_type_id.String())
	}
	result = arr[0]._number
	for ; arrIndex < arrSize; arrIndex++ {
		if !arr[arrIndex].IsNumber() {
			return CreateNil(), errors.New("need arguments type is number, but got " + arr[arrIndex]._sexp_type_id.String())
		}
		result = result - arr[arrIndex]._number
	}
	return CreateInt(result), nil
}

func NewSub() *Sexpression {
	return CreateSubroutine("sub", _sub_minus_Apply)
}

func _sub_mul_Apply(self *Sexpression, ctx context.Context, _ *Sexpression, args *Sexpression) (*Sexpression, error) {
	if !args.IsConsCell() {
		return CreateNil(), errors.New("need arguments")
	}

	if IsEmptyList(args) {
		return CreateInt(0), nil
	}

	arr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}
	arrIndex := uint64(1)

	var result = int64(1)
	if !arr[0].IsNumber() {
		return CreateNil(), errors.New("need arguments type is number, but got " + arr[0]._sexp_type_id.String())
	}
	result = arr[0]._number
	for ; arrIndex < arrSize; arrIndex++ {
		if !arr[arrIndex].IsNumber() {
			return CreateNil(), errors.New("need arguments type is number, but got " + arr[arrIndex]._sexp_type_id.String())
		}
		result = result * arr[arrIndex]._number
	}
	return CreateInt(result), nil
}

func NewMul() *Sexpression {
	return CreateSubroutine("*", _sub_mul_Apply)
}

func _sub_div_Apply(self *Sexpression, ctx context.Context, _ *Sexpression, args *Sexpression) (*Sexpression, error) {
	if !args.IsConsCell() {
		return CreateNil(), errors.New("need arguments")
	}

	if IsEmptyList(args) {
		return CreateInt(0), nil
	}

	arr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}
	arrIndex := uint64(1)

	var result = int64(1)
	if !arr[0].IsNumber() {
		return CreateNil(), errors.New("need arguments type is number, but got " + arr[0]._sexp_type_id.String())
	}
	result = arr[0]._number
	for ; arrIndex < arrSize; arrIndex++ {
		if !arr[arrIndex].IsNumber() {
			return CreateNil(), errors.New("need arguments type is number, but got " + arr[arrIndex]._sexp_type_id.String())
		}
		result = result / arr[arrIndex]._number
	}
	return CreateInt(result), nil
}

func NewDiv() *Sexpression {
	return CreateSubroutine("/", _sub_div_Apply)
}

func _sub_mod_Apply(self *Sexpression, ctx context.Context, _ *Sexpression, args *Sexpression) (*Sexpression, error) {
	if !args.IsConsCell() {
		return CreateNil(), errors.New("need arguments")
	}

	if IsEmptyList(args) {
		return CreateInt(0), nil
	}

	arr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}
	arrIndex := uint64(1)

	var result = int64(1)
	if !arr[0].IsNumber() {
		return CreateNil(), errors.New("need arguments type is number, but got " + arr[0]._sexp_type_id.String())
	}
	result = arr[0]._number
	for ; arrIndex < arrSize; arrIndex++ {
		if !arr[arrIndex].IsNumber() {
			return CreateNil(), errors.New("need arguments type is number, but got " + arr[arrIndex]._sexp_type_id.String())
		}
		result = result % arr[arrIndex]._number
	}
	return CreateInt(result), nil
}

func NewMod() *Sexpression {
	return CreateSubroutine("%", _sub_mod_Apply)
}
