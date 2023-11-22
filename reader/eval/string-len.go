package eval

import (
	"context"
	"errors"
	"unicode/utf8"
)

func _subr_string_len_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)
	if err != nil {
		return CreateNil(), err
	}
	if arrSize != 1 {
		return CreateNil(), errors.New("need args size is 1")
	}

	size := int64(utf8.RuneCountInString(arr[0]._string))

	return CreateInt(size), nil
}

func NewStringLen() *Sexpression {
	return CreateSubroutine("string-len", _subr_string_len_Apply)
}
