package eval

import "context"

type _or struct {
}

func (_ _or) TypeId() string {
	return "special_form.or"
}

func (_ _or) AtomId() AtomType {
	return AtomTypeSpecialForm
}

func (_ _or) String() string {
	return "#<syntax #or>"
}

func (_ _or) IsList() bool {
	return false
}

func (a _or) Equals(sexp SExpression) bool {
	return a.TypeId() == sexp.TypeId()
}

func (_ _or) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	evaluatedElm := NewConsCell(NewNil(), NewNil()).(SExpression)
	var err error

	for i := uint64(0); i < argsLength; i++ {
		evaluatedElm, err = Eval(ctx, args[i], env)
		if err != nil {
			return nil, err
		}
		if !NewBool(false).Equals(evaluatedElm) {
			return evaluatedElm, nil
		}
	}

	return evaluatedElm, nil
}

func NewOr() SExpression {
	return &_or{}
}
