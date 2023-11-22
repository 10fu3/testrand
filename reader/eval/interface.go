package eval

import "context"

type Callable interface {
	SExpression
	Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error)
}
