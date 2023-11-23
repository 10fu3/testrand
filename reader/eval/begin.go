package eval

import "context"

type _begin struct{}

func (_ *_begin) TypeId() string {
	return "special_form.begin"
}

func (_ *_begin) AtomId() AtomType {
	return AtomTypeSpecialForm
}

func (_ *_begin) String() string {
	return "#<syntax #begin>"
}

func (_ *_begin) IsList() bool {
	return false
}

func (b *_begin) Equals(sexp SExpression) bool {
	return b.TypeId() == sexp.TypeId()
}

func (_ *_begin) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	last := NewNil()
	var err error

	for _, sexp := range args {
		last, err = Eval(ctx, sexp, env)
		if err != nil {
			return nil, err
		}
	}
	return last, err
}

func NewBegin() SExpression {
	return &_begin{}
}
