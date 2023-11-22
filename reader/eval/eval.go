package eval

import (
	"context"
	"errors"
)

func Eval(ctx context.Context, sexp *Sexpression, env *Sexpression) (*Sexpression, error) {
	sexpType := sexp.SexpressionTypeId()
	switch sexpType {
	case SexpressionTypeNumber,
		SexpressionTypeString,
		SexpressionTypeBool,
		SexpressionTypeEnvironment,
		SexpressionTypeNil,
		SexpressionTypeNativeArray,
		SexpressionTypeNativeHashmap:
		return sexp, nil
	case SexpressionTypeSymbol:
		v, ok := env.GetValueFromFrame(sexp)
		if !ok {
			return CreateNil(), errors.New("symbol not found: " + sexp.String())
		}
		return v, nil
	case SexpressionTypeConsCell:
		cell := sexp._cell
		applied, err := Eval(ctx, cell._car, env)
		if err != nil {
			return CreateNil(), err
		}
		appliedType := applied.SexpressionTypeId()
		if err != nil {
			return CreateNil(), err
		}
		if SexpressionTypeClosure == appliedType || SexpressionTypeSubroutine == appliedType {
			appliedArgs, appliedArgsErr := evalArgument(ctx, cell._cdr, env)
			if appliedArgsErr != nil {
				return CreateNil(), err
			}
			return applied._applyFunc(applied, ctx, env, appliedArgs)
		}

		if SexpressionTypeClosure == appliedType || SexpressionTypeSubroutine == appliedType {
			if !sexp.IsList() {
				return CreateNil(), errors.New("not list")
			}
			arr, arrSize, toArrErr := ToArray(cell._cdr)
			if toArrErr != nil {
				return CreateNil(), toArrErr
			}

			rootEvaluetedArgs := CreateEmptyList()
			lookEvalutedArg := rootEvaluetedArgs
			for i, v := range arr {
				evaluatedArg, evalutedArgErr := Eval(ctx, v, env)
				if evalutedArgErr != nil {
					return CreateNil(), err
				}
				lookEvalutedArg._cell._car = evaluatedArg
				lookEvalutedArg._cell._cdr = CreateEmptyList()
				if uint64(i) < arrSize {
					lookEvalutedArg = lookEvalutedArg._cell._cdr
				}
			}
			if err != nil {
				return CreateNil(), err
			}
			return applied._applyFunc(applied, ctx, env, rootEvaluetedArgs)
		}

		if SexpressionTypeSpecialForm == appliedType {
			return applied._applyFunc(applied, ctx, env, cell._cdr)
		}

	}
	return CreateNil(), errors.New("unknown eval: " + sexp.String())
}

func evalArgument(ctx context.Context, sexp *Sexpression, env *Sexpression) (*Sexpression, error) {
	if SexpressionTypeConsCell != sexp._sexp_type_id {
		return Eval(ctx, sexp, env)
	}

	cell := sexp._cell

	carEvaluated, err := Eval(ctx, cell._car, env)
	if err != nil {
		return CreateNil(), err
	}

	cdrEvaluated, err := evalArgument(ctx, cell._cdr, env)
	if err != nil {
		return CreateNil(), err
	}

	return CreateConsCell(carEvaluated, cdrEvaluated), nil
}

func _subr_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	argsArr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}

	if arrSize != 2 {
		return CreateNil(), errors.New("malformed eval")
	}

	targetEnv, err := Eval(ctx, argsArr[1], env)

	if err != nil {
		return CreateNil(), err
	}

	return Eval(ctx, argsArr[0], targetEnv)
}

func NewEval() *Sexpression {
	return CreateSubroutine("eval", _subr_Apply)
}
