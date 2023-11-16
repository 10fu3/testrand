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
		SExpressionTypeNil:
		return sexp, nil
	case SExpressionTypeSymbol:
		if v, _ := env.GetValue(sexp.(Symbol)); v == nil {
			if ctx.Value("transaction") == nil {
				return nil, errors.New("unknown symbol")
			}

		}
		return env.GetValue(sexp.(Symbol))
	case SExpressionTypeConsCell:
		cell := sexp.(ConsCell)
		applied, err := Eval(ctx, cell.GetCar(), env)
		appliedType := applied.SExpressionTypeId()
		if err != nil {
			return nil, err
		}
		if SExpressionTypeClosure == appliedType || SExpressionTypeSubroutine == appliedType {
			appliedArgs, err := evalArgument(ctx, cell.GetCdr(), env)
			if err != nil {
				return nil, err
			}
			return applied.(Callable).Apply(ctx, env, appliedArgs)
		}
		if SExpressionTypeSpecialForm == appliedType {
			if err != nil {
				return nil, err
			}
			return applied.(Callable).Apply(ctx, env, cell.GetCdr())
		}

	}
	return nil, errors.New("unknown eval: " + sexp.String())
}

func evalArgument(ctx context.Context, sexp SExpression, env Environment) (SExpression, error) {
	if "cons_cell" != sexp.TypeId() {
		return Eval(ctx, sexp, env)
	}

	cell := sexp.(ConsCell)

	carEvaluated, err := Eval(ctx, cell.GetCar(), env)
	if err != nil {
		return nil, err
	}

	cdrEvaluated, err := evalArgument(ctx, cell.GetCdr(), env)
	if err != nil {
		return nil, err
	}

	return NewConsCell(carEvaluated, cdrEvaluated), nil
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

func (_ *_eval) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	argsArr, err := ToArray(args)

	if err != nil {
		return nil, err
	}

	if len(argsArr) != 2 {
		return nil, errors.New("malformed eval")
	}

	targetEnv, err := Eval(ctx, argsArr[1].(Environment), env)

	if err != nil {
		return nil, err
	}

	return Eval(ctx, argsArr[0], targetEnv.(Environment))
}

func NewEval() SExpression {
	return &_eval{}
}
