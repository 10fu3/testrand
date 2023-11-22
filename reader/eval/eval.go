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
			if !cell.IsList() {
				return nil, errors.New("not list")
			}
			arr, toArrErr := ToArray(cell.GetCdr())
			if toArrErr != nil {
				return nil, toArrErr
			}

			rootEvaluetedArgs := &_cons_cell{}
			lookEvalutedArg := rootEvaluetedArgs
			for i, v := range arr {
				evaluatedArg, evalutedArgErr := Eval(ctx, v, env)
				if evalutedArgErr != nil {
					return nil, err
				}
				if evaluatedArg == nil {
					return nil, errors.New("evaluated arg is nil")
				}
				lookEvalutedArg.Car = evaluatedArg
				lookEvalutedArg.Cdr = NewConsCell(NewNil(), NewNil())
				if i < len(arr) {
					lookEvalutedArg = (lookEvalutedArg.Cdr).(*_cons_cell)
				}
			}
			if err != nil {
				return nil, err
			}
			return applied.(Callable).Apply(ctx, env, rootEvaluetedArgs)
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
