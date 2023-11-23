package eval

import "context"

type _and struct {
}

func (_ _and) TypeId() string {
	return "special_form.and"
}

func (_ _and) AtomId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ _and) String() string {
	return "#<syntax #and>"
}

func (_ _and) IsList() bool {
	return false
}

func (a _and) Equals(sexp SExpression) bool {
	return a.TypeId() == sexp.TypeId()
}

func (_ _and) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	evaluatedElm := NewConsCell(NewNil(), NewNil()).(SExpression)

	for i := uint64(0); i < argsLength; i++ {
		evaluatedElm, err := Eval(ctx, args[i], env)
		if err != nil {
			return nil, err
		}
		if NewBool(false).Equals(evaluatedElm) {
			return evaluatedElm, nil
		}
	}

	return evaluatedElm, nil
}

func NewAnd() SExpression {
	return &_and{}
}
