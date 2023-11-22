package eval

import "context"

func _sub_And(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {

	if IsEmptyList(args) {
		return CreateBool(true), nil
	}

	arr, arrSize, err := ToArray(args)

	if err != nil {
		return &Sexpression{}, err
	}

	evaluatedElm := CreateEmptyList()

	for i := uint64(0); i < arrSize; i++ {
		evaluatedElm, err = Eval(ctx, arr[i], env)
		if err != nil {
			return CreateNil(), err
		}
		if CreateBool(false).Equals(evaluatedElm) {
			return evaluatedElm, nil
		}
	}

	return evaluatedElm, nil
}

func NewAnd() *Sexpression {
	return CreateSubroutine("and", _sub_And)
}
