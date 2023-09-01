package eval

import "context"

type Callable interface {
	SExpression
	Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error)
}
