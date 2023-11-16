package eval

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type _closure struct {
	id           string
	body         SExpression
	formals      SExpression
	env          Environment
	formalsCount int
}

type Closure interface {
	SExpression
	GetFormalsCount() int
}

func NewClosure(body SExpression, formals SExpression, env Environment, formalsCount int) Callable {
	return &_closure{body: body, formals: formals, env: env, id: uuid.NewString(), formalsCount: formalsCount}
}

func (c *_closure) TypeId() string {
	return "closure"
}

func (c *_closure) SExpressionTypeId() SExpressionType {
	return SExpressionTypeClosure
}

func (c *_closure) String() string {
	return fmt.Sprintf("#<closure %s>", c.id)
}

func (c *_closure) IsList() bool {
	return false
}

func (c *_closure) Equals(args SExpression) bool {
	if "closure" != args.TypeId() {
		return false
	}
	return c.id == args.(*_closure).id
}

func (c *_closure) GetFormalsCount() int {
	return c.formalsCount
}

func (c *_closure) Apply(ctx context.Context, _ Environment, args SExpression) (SExpression, error) {
	loopFormals := c.formals
	loopArgs := args
	env, _ := NewEnvironment(c.env)

	for {
		if IsEmptyList(loopFormals) {
			if IsEmptyList(loopArgs) {
				break
			}
			return nil, errors.New("argument size more than formals")
		}
		if "symbol" == loopFormals.TypeId() {
			env.Define(loopFormals.(Symbol), loopArgs)
			break
		}
		if "cons_cell" == loopFormals.TypeId() {
			cellFormals := loopFormals.(ConsCell)
			if "symbol" != cellFormals.GetCar().TypeId() {
				return nil, errors.New("need symbol")
			}
			if "cons_cell" != loopArgs.TypeId() {
				return nil, errors.New("argument size less than formals")
			}
			cellArgs := loopArgs.(ConsCell)
			env.Define(cellFormals.GetCar().(Symbol), cellArgs.GetCar())
			loopFormals = cellFormals.GetCdr()
			loopArgs = cellArgs.GetCdr()
		}
	}
	return Eval(ctx, c.body, env)
}
