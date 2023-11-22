package eval

import "context"

func _syntax_or_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {

	if IsEmptyList(args) {
		return CreateBool(false), nil
	}

	arr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}

	evaluatedElm := CreateEmptyList()

	for i := uint64(0); i < arrSize; i++ {
		evaluatedElm, err = Eval(ctx, arr[i], env)
		if err != nil {
			return CreateNil(), err
		}
		if !CreateBool(false).Equals(evaluatedElm) {
			return evaluatedElm, nil
		}
	}

	return evaluatedElm, nil
}

func NewOr() *Sexpression {
	return CreateSpecialForm("or", _syntax_or_Apply)
}
