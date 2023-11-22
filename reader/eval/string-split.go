package eval

import (
	"context"
	"errors"
	"strings"
)

func _subr_string_split_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)
	if err != nil {
		return CreateNil(), err
	}
	if arrSize < 1 {
		return CreateNil(), errors.New("need args size is 1")
	}

	if !arr[0].IsString() {
		return CreateNil(), errors.New("need args type is string")
	}

	str := arr[0]._string

	var sep string
	if arrSize == 2 {
		if !arr[1].IsString() {
			return CreateNil(), errors.New("need args type is string")
		}
		sep = arr[1]._string
	} else {
		sep = ""
	}

	split := strings.Split(str, sep)
	converted := make([]*Sexpression, len(split))
	for i := 0; i < len(split); i++ {
		converted[i] = CreateString(split[i])
	}

	return CreateNativeArray(converted), nil
}

func NewStringSplit() *Sexpression {
	return CreateSubroutine("string-split", _subr_string_split_Apply)
}
