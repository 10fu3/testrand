package eval

import (
	"context"
	"errors"
)

func Eval(ctx context.Context, sexp SExpression, env Environment) (SExpression, error) {
	sexpType := sexp.AtomId()
	switch sexpType {
	case AtomTypeNumber,
		AtomTypeString,
		AtomTypeBool,
		AtomTypeEnvironment,
		AtomTypeNil,
		AtomTypeNativeArray,
		AtomTypeNativeHashmap:
		return sexp, nil
	case AtomTypeSymbol:
		v, err := env.GetValue(sexp.(Symbol))

		if err != nil {
			return nil, err
		}
		return v, nil
	case AtomTypeConsCell:
		cell := sexp.(ConsCell)
		applied, err := Eval(ctx, cell.GetCar(), env)
		if err != nil {
			return nil, err
		}
		appliedType := applied.AtomId()
		if err != nil {
			return nil, err
		}
		if AtomTypeClosure == appliedType || AtomTypeSubroutine == appliedType {
			argsArr, argsArrSize, argsErr := ToArray(cell.GetCdr())
			appliedArgs := make([]SExpression, argsArrSize)

			for i := uint64(0); i < argsArrSize && argsErr == nil; i++ {
				appliedArgs[i], argsErr = Eval(ctx, argsArr[i], env)
			}

			if argsErr != nil {
				return nil, argsErr
			}

			return applied.(Callable).Apply(ctx, env, appliedArgs, argsArrSize)
		}
		if AtomTypeSpecialForm == appliedType {
			args, length, toArrErr := ToArray(cell.GetCdr())
			if toArrErr != nil {
				return nil, err
			}
			return applied.(Callable).Apply(ctx, env, args, length)
		}

	}
	return nil, errors.New("unknown eval: " + sexp.String())
}

type _eval struct{}

func (_ *_eval) TypeId() string {
	return "subroutine.eval"
}

func (_ *_eval) AtomId() AtomType {
	return AtomTypeSubroutine
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
