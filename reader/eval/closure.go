package eval

import (
	"context"
	"errors"
	"fmt"
	"testrand/cmap"
)

func ____max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func _closure_run(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {

	loopArgs, loopArgsSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}

	if loopArgsSize != self._apply_formalsCount {
		return CreateNil(), errors.New(fmt.Sprintf("not match argument size: required %d but got %d", self._apply_formalsCount, loopArgsSize))
	}

	frame := cmap.New[*Sexpression]()

	var argElem *Sexpression = CreateNil()

	for formalsIndex, formalElem := range self._apply_formals {
		argElem = loopArgs[formalsIndex]
		if formalElem.IsSymbol() {
			frame.Set(formalElem._symbol._string, argElem)
			break
		}
		if formalElem.IsConsCell() {
			cellFormals := formalElem._cell
			if SexpressionTypeSymbol != cellFormals._car._sexp_type_id {
				return CreateNil(), errors.New("need symbol")
			}
			if SexpressionTypeSymbol != argElem._sexp_type_id {
				return CreateNil(), errors.New("argument size less than formals")
			}
			cellArgs := argElem._cell
			frame.Set(cellFormals._car._symbol._string, cellArgs._car)
		}
	}
	newEnv, err := NewEnvironmentForClosure(self._apply_env, frame)

	if err != nil {
		return CreateNil(), err
	}
	return Eval(ctx, self._apply_body, newEnv)
}
