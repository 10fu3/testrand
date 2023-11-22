package eval

import (
	"context"
	"errors"
)

type _quasiquote struct {
}

func _innerEvalQuasiquote(self *Sexpression, ctx context.Context, env *Sexpression, x *Sexpression) (*Sexpression, error) {

	if SexpressionTypeConsCell != x._sexp_type_id {
		return CreateNil(), errors.New("malformed quote")
	}

	pair := x._cell
	car := pair._car
	cdr := pair._cdr
	if car.Equals(CreateSymbol("unquote")) {
		if SexpressionTypeConsCell != cdr._sexp_type_id {
			return CreateNil(), errors.New("unquote must be followed by a list")
		}
		unquoted, err := Eval(ctx, cdr._cell._car, env)
		return unquoted, err
	}

	if car.Equals(CreateSymbol("quasiquote")) {
		return x, nil
	}

	if SexpressionTypeConsCell == car._sexp_type_id && (car._cell._car).Equals(CreateSymbol("unquote-splicing")) {
		innerPair := car._cell._cdr._cell
		innerPairCarQuoteEvaluated, err := _innerEvalQuasiquote(self, ctx, env, innerPair._car)
		if err != nil {
			return CreateNil(), err
		}
		innerPairCarEvalueted, err := Eval(ctx, innerPairCarQuoteEvaluated, env)
		if err != nil {
			return CreateNil(), err
		}
		if !innerPairCarEvalueted.IsList() {
			return CreateNil(), errors.New("unquote-splicing must be followed by a list")
		}
		if IsEmptyList(innerPair._cdr) {
			cdrEvaluated, err := _innerEvalQuasiquote(self, ctx, env, cdr)
			if err != nil {
				return CreateNil(), err
			}
			joined, err := JoinList(innerPairCarEvalueted, cdrEvaluated)

			if err != nil {
				return CreateNil(), err
			}

			return joined, nil
		}
		innerPairCdrEvaluated, err := _innerEvalQuasiquote(self, ctx, env, innerPair._cdr)
		if err != nil {
			return CreateNil(), err
		}
		return CreateConsCell(innerPairCarEvalueted, innerPairCdrEvaluated), nil
	}
	carEvaluated, err := _innerEvalQuasiquote(self, ctx, env, car)
	if err != nil {
		return CreateNil(), err
	}
	if IsEmptyList(cdr) {
		return CreateConsCell(carEvaluated, CreateConsCell(CreateNil(), CreateNil())), nil
	}

	cdrEvaluated, err := _innerEvalQuasiquote(self, ctx, env, cdr)
	if err != nil {
		return CreateNil(), err
	}

	return CreateConsCell(carEvaluated, cdrEvaluated), nil
}

// this function is lisp interpter function for quasiquote
func _syntax_quasiquote_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}
	if arrSize != 1 {
		return CreateNil(), errors.New("malformed quote")
	}
	return _innerEvalQuasiquote(self, ctx, env, arr[0])
}

func NewQuasiquote() *Sexpression {
	return CreateSpecialForm("quasiquote", _syntax_quasiquote_Apply)
}
