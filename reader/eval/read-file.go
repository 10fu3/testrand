package eval

import (
	"context"
	"errors"
	"os"
)

func _subr_read_file_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {

	args, argsLen, err := ToArray(arguments)
	if err != nil {
		return CreateNil(), err
	}

	if 1 > argsLen {
		return CreateNil(), errors.New("need arguments size is 1")
	}

	rawPath := args[0]
	if !rawPath.IsString() {
		return CreateNil(), errors.New("need arguments type is string")
	}

	path := rawPath._string

	fileReader, err := os.ReadFile(path)

	if err != nil {
		return CreateNil(), err
	}

	return CreateString(string(fileReader)), nil
}

func NewReadFile() *Sexpression {
	return CreateSubroutine("read-file", _subr_read_file_Apply)
}
