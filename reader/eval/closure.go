package eval

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type _closure struct {
	id      string
	body    SExpression
	formals SExpression
	env     Environment
}

func NewClosure(body SExpression, formals SExpression, env Environment) Callable {
	return &_closure{body: body, formals: formals, env: env, id: uuid.NewString()}
}

func (c *_closure) Type() string {
	return "callable.closure"
}

func (c *_closure) String() string {
	return fmt.Sprintf("#<closure %s>", c.id)
}

func (c *_closure) IsList() bool {
	return false
}

func (c *_closure) Equals(args SExpression) bool {
	if "callable.closure" != args.Type() {
		return false
	}
	return c.id == args.(*_closure).id
}

func (c *_closure) Apply(_ Environment, args SExpression) (SExpression, error) {
	loopFormals := c.formals
	loopArgs := args
	env := NewEnvironment(c.env)

	for {
		if "cons_cell" == loopFormals.Type() {
			f := loopFormals.(ConsCell)
			if "nil" == f.GetCar().Type() && "nil" == f.GetCdr().Type() {
				if "cons_cell" == loopArgs.Type() {
					a := loopArgs.(ConsCell)
					if "nil" == a.GetCar().Type() && "nil" == a.GetCdr().Type() {
						break
					}
				}
				return nil, errors.New("argument size more than formals")
			}

			if "symbol" != f.GetCar().Type() {
				return nil, errors.New("type error: " + f.GetCar().String())
			}
			if "cons_cell" != loopArgs.Type() {
				return nil, errors.New("argument size less than formals")
			}
			cellArgs := loopArgs.(ConsCell)
			env.Define(f.GetCar().(Symbol), cellArgs.GetCar())
			loopFormals = f.GetCdr()
			loopArgs = cellArgs.GetCdr()
		}
		if "symbol" == loopFormals.Type() {
			env.Define(loopFormals.(Symbol), loopArgs)
			break
		}
	}
	return Eval(c.body, env)
}
