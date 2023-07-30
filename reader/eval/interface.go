package eval

type Callable interface {
	SExpression
	Apply(env Environment, args SExpression) (SExpression, error)
}
