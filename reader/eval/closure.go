package eval

import (
	"context"
	"errors"
	"fmt"
)

type _closure struct {
	body         SExpression
	formals      []SExpression
	env          Environment
	formalsCount int
}

type Closure interface {
	SExpression
	GetFormalsCount() int
	Callable
}

func NewClosure(body SExpression, formals []SExpression, env Environment, formalsCount int) (Callable, error) {
	return &_closure{body: body, formals: formals, env: env, formalsCount: formalsCount}, nil
}

func (c *_closure) TypeId() string {
	return "closure"
}

func (c *_closure) SExpressionTypeId() SExpressionType {
	return SExpressionTypeClosure
}

func (c *_closure) String() string {
	return fmt.Sprintf("#<closure>")
}

func (c *_closure) IsList() bool {
	return false
}

func (c *_closure) Equals(args SExpression) bool {
	if "closure" != args.TypeId() {
		return false
	}
	return args.(*_closure) == c
}

func (c *_closure) GetFormalsCount() int {
	return c.formalsCount
}

func (c *_closure) Apply(ctx context.Context, _ Environment, args SExpression) (SExpression, error) {

	loopArgs, err := ToArray(args)

	if err != nil {
		return nil, err
	}

	if len(loopArgs) != len(c.formals) {
		return nil, errors.New(fmt.Sprintf("not match argument size: %d != %d", len(loopArgs), len(c.formals)))
	}

	frame := map[string]SExpression{}

	var argElem SExpression = NewNil()

	for formalsIndex, formalElem := range c.formals {
		argElem = loopArgs[formalsIndex]
		if "symbol" == formalElem.TypeId() {
			frame[formalElem.(Symbol).GetValue()] = argElem
			break
		}
		if "cons_cell" == formalElem.TypeId() {
			cellFormals := formalElem.(ConsCell)
			if "symbol" != cellFormals.GetCar().TypeId() {
				return nil, errors.New("need symbol")
			}
			if "cons_cell" != argElem.TypeId() {
				return nil, errors.New("argument size less than formals")
			}
			cellArgs := argElem.(ConsCell)
			frame[cellFormals.GetCar().(Symbol).GetValue()] = cellArgs.GetCar()
		}
	}
	env, err := NewEnvironmentForClosure(c.env, frame)

	if err != nil {
		return nil, err
	}
	return Eval(ctx, c.body, env)
}
