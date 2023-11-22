package eval

import (
	"context"
	"errors"
)

func Eval(ctx context.Context, sexp SExpression, env Environment) (SExpression, error) {
	sexpType := sexp.SExpressionTypeId()
	switch sexpType {
	case SExpressionTypeNumber,
		SExpressionTypeString,
		SExpressionTypeBool,
		SExpressionTypeEnvironment,
		SExpressionTypeNil,
		SExpressionTypeNativeArray,
		SExpressionTypeNativeHashmap:
		return sexp, nil
	case SExpressionTypeSymbol:

		v, err := env.GetValue(sexp.(Symbol))

		if err != nil {
			return nil, err
		}
		return v, nil
	case SExpressionTypeConsCell:
		cell := sexp.(ConsCell)
		applied, err := Eval(ctx, cell.GetCar(), env)
		if err != nil {
			return nil, err
		}
		appliedType := applied.SExpressionTypeId()
		if err != nil {
			return nil, err
		}
		if SExpressionTypeClosure == appliedType || SExpressionTypeSubroutine == appliedType {
			appliedArgs, argsSize, err := evalArgument(ctx, cell.GetCdr(), env)
			if err != nil {
				return nil, err
			}
			for i, j := 0, len(appliedArgs)-1; i < j; i, j = i+1, j-1 {
				appliedArgs[i], appliedArgs[j] = appliedArgs[j], appliedArgs[i]
			}
			return applied.(Callable).Apply(ctx, env, appliedArgs, argsSize)
		}
		if SExpressionTypeSpecialForm == appliedType {
			args, length, toArrErr := ToArray(cell.GetCdr())
			if toArrErr != nil {
				return nil, err
			}
			return applied.(Callable).Apply(ctx, env, args, length)
		}

	}
	return nil, errors.New("unknown eval: " + sexp.String())
}

func evalArgument(ctx context.Context, sexp SExpression, env Environment) ([]SExpression, uint64, error) {
	if "cons_cell" != sexp.TypeId() {
		result, err := Eval(ctx, sexp, env)
		return []SExpression{result}, 1, err
	}

	if IsEmptyList(sexp) {
		return []SExpression{}, 0, nil
	}

	cell := sexp.(ConsCell)

	carEvaluated, err := Eval(ctx, cell.GetCar(), env)
	if err != nil {
		return nil, 0, err
	}

	cdrEvaluated, size, err := evalArgument(ctx, cell.GetCdr(), env)
	if err != nil {
		return nil, 0, err
	}

	if len(cdrEvaluated)+1 < cap(cdrEvaluated) {
		cdrEvaluated = cdrEvaluated[:len(cdrEvaluated)+1] // slice の延長
		cdrEvaluated[len(cdrEvaluated)] = carEvaluated
	} else if cap(cdrEvaluated) < len(cdrEvaluated)+1 {
		cdrEvaluated = append(cdrEvaluated, carEvaluated)
	}

	return cdrEvaluated, size + 1, nil
}

type _eval struct{}

func (_ *_eval) TypeId() string {
	return "subroutine.eval"
}

func (_ *_eval) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_eval) String() string {
	return "#<subr eval>"
}

func (_ *_eval) IsList() bool {
	return false
}

func (e *_eval) Equals(sexp SExpression) bool {
	return e.TypeId() == sexp.TypeId()
}

func (_ *_eval) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("malformed eval")
	}

	targetEnv, err := Eval(ctx, args[1].(Environment), env)

	if err != nil {
		return nil, err
	}

	return Eval(ctx, args[0], targetEnv.(Environment))
}

func NewEval() SExpression {
	return &_eval{}
}
