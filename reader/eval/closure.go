package eval

import (
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

func (c *_closure) Type() string {
	return "closure"
}

func (c *_closure) String() string {
	return fmt.Sprintf("#<closure %s>", c.id)
}

func (c *_closure) IsList() bool {
	return false
}

func (c *_closure) Equals(args SExpression) bool {
	if "closure" != args.Type() {
		return false
	}
	return c.id == args.(*_closure).id
}

func (c *_closure) GetFormalsCount() int {
	return c.formalsCount
}

func (c *_closure) Apply(_ Environment, args SExpression) (SExpression, error) {
	loopFormals := c.formals
	loopArgs := args
	env := NewEnvironment(c.env)

	for {
		if IsEmptyList(loopFormals) {
			if IsEmptyList(loopArgs) {
				break
			}
			return nil, errors.New("argument size more than formals")
		}
		if "symbol" == loopFormals.Type() {
			env.Define(loopFormals.(Symbol), loopArgs)
			break
		}
		if "cons_cell" == loopFormals.Type() {
			cellFormals := loopFormals.(ConsCell)
			if "symbol" != cellFormals.GetCar().Type() {
				return nil, errors.New("need symbol")
			}
			if "cons_cell" != loopArgs.Type() {
				return nil, errors.New("argument size less than formals")
			}
			cellArgs := loopArgs.(ConsCell)
			env.Define(cellFormals.GetCar().(Symbol), cellArgs.GetCar())
			loopFormals = cellFormals.GetCdr()
			loopArgs = cellArgs.GetCdr()
		}
	}
	return Eval(c.body, env)
}
