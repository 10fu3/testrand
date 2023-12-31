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
	formalsCount uint64
}

type Closure interface {
	SExpression
	GetFormalsCount() uint64
	Callable
}

func NewClosure(body SExpression, formals []SExpression, env Environment, formalsCount uint64) (Callable, error) {
	return &_closure{body: body, formals: formals, env: env, formalsCount: formalsCount}, nil
}

func (c *_closure) TypeId() string {
	return "closure"
}

func (c *_closure) AtomId() AtomType {
	return AtomTypeClosure
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

func (c *_closure) GetFormalsCount() uint64 {
	return c.formalsCount
}

func (c *_closure) Apply(ctx context.Context, _ Environment, loopArgs []SExpression, loopArgsLength uint64) (SExpression, error) {

	if loopArgsLength != c.formalsCount {
		return nil, errors.New(fmt.Sprintf("not match argument size: %d != %d", len(loopArgs), len(c.formals)))
	}

	frame := make(map[string]SExpression, c.formalsCount)

	var argElem SExpression = NewNil()

	for formalsIndex := uint64(0); formalsIndex < c.formalsCount; formalsIndex++ {
		argElem = loopArgs[formalsIndex]
		formalsElem := c.formals[formalsIndex]
		if AtomTypeSymbol == formalsElem.AtomId() {
			frame[formalsElem.(Symbol).GetValue()] = loopArgs[formalsIndex]
			break
		}
		if AtomTypeConsCell == formalsElem.AtomId() {
			cellFormals := formalsElem.(ConsCell)
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
