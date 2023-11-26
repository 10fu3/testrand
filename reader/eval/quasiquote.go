package eval

import (
	"context"
	"errors"
)

type _quasiquote struct {
}

func (_ *_quasiquote) TypeId() string {
	return "special_form.quasiquote"
}

func (_ *_quasiquote) AtomId() AtomType {
	return AtomTypeSpecialForm
}

func (_ *_quasiquote) String() string {
	return "#<syntax #quasiquote>"
}

func (_ *_quasiquote) IsList() bool {
	return false
}

func (q *_quasiquote) Equals(sexp SExpression) bool {
	return q.TypeId() == sexp.TypeId()
}

func _innerEvalQuasiquote(ctx context.Context, env Environment, x SExpression) (SExpression, error) {
	if x.TypeId() != "cons_cell" {
		return x, nil
	}
	pair := x.(ConsCell)
	car := pair.GetCar()
	cdr := pair.GetCdr()
	if car.Equals(NewSymbol("unquote")) {
		if cdr.TypeId() != "cons_cell" {
			return nil, errors.New("unquote must be followed by a list")
		}
		unquoted, err := Eval(ctx, cdr.(ConsCell).GetCar(), env)
		return unquoted, err
	}

	if car.Equals(NewSymbol("quasiquote")) {
		return x, nil
	}

	if car.TypeId() == "cons_cell" && (car.(ConsCell).GetCar()).Equals(NewSymbol("unquote-splicing")) {
		innerPair := car.(ConsCell).GetCdr().(ConsCell)
		innerPairCarQuoteEvaluated, err := _innerEvalQuasiquote(ctx, env, innerPair.GetCar())
		if err != nil {
			return nil, err
		}
		innerPairCarEvaluated, err := Eval(ctx, innerPairCarQuoteEvaluated, env)
		if err != nil {
			return nil, err
		}
		if !innerPairCarEvaluated.IsList() {
			return nil, errors.New("unquote-splicing must be followed by a list")
		}
		if IsEmptyList(innerPair.GetCdr()) {
			cdrEvaluated, cdrEvaluatedErr := _innerEvalQuasiquote(ctx, env, cdr)
			if cdrEvaluatedErr != nil {
				return nil, cdrEvaluatedErr
			}
			joined, cdrEvaluatedErr := JoinList(innerPairCarEvaluated, cdrEvaluated)

			if cdrEvaluatedErr != nil {
				return nil, cdrEvaluatedErr
			}

			return joined, nil
		}
		innerPairCdrEvaluated, err := _innerEvalQuasiquote(ctx, env, innerPair.GetCdr())
		if err != nil {
			return nil, err
		}
		return NewConsCell(innerPairCarEvaluated, innerPairCdrEvaluated), nil
	}
	carEvaluated, err := _innerEvalQuasiquote(ctx, env, car)
	if err != nil {
		return nil, err
	}
	if IsEmptyList(cdr) {
		return NewConsCell(carEvaluated, NewConsCell(NewNil(), NewNil())), nil
	}

	cdrEvaluated, err := _innerEvalQuasiquote(ctx, env, cdr)
	if err != nil {
		return nil, err
	}

	return NewConsCell(carEvaluated, cdrEvaluated), nil
}

// this function is lisp interpter function for quasiquote
func (_ *_quasiquote) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	if argsLength != 1 {
		return nil, errors.New("malformed quote")
	}
	return _innerEvalQuasiquote(ctx, env, args[0])
}

func NewQuasiquote() SExpression {
	return &_quasiquote{}
}
