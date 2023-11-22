package eval

import "context"

type _begin struct{}

func _begin_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, _, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}

	last := CreateNil()

	for _, sexp := range arr {
		last, err = Eval(ctx, sexp, env)
		if err != nil {
			return CreateNil(), err
		}
	}
	return last, err
}

func NewBegin() *Sexpression {
	return CreateSpecialForm("begin", _begin_Apply)
}
