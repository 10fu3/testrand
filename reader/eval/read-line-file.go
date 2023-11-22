package eval

import (
	"bufio"
	"context"
	"errors"
	"os"
)

func _subr_file_read_line_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	// 1st filepath string
	// 2nd on-load-line function
	// 3rd on-load-end function
	args, argsLen, err := ToArray(arguments)
	if err != nil {
		return CreateNil(), err
	}

	if 2 > argsLen {
		return CreateNil(), errors.New("need arguments size is 1")
	}

	rawPath := args[0]
	if !rawPath.IsString() {
		return CreateNil(), errors.New("need arguments type is string")
	}

	path := rawPath._string

	onLoadLine := args[1]

	var fp *os.File

	fp, err = os.Open(path)
	if err != nil {
		return CreateNil(), err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		evalTarget := CreateConsCell(onLoadLine,
			CreateConsCell(CreateString(scanner.Text()),
				CreateConsCell(CreateNil(), CreateNil())))

		Eval(ctx, evalTarget, env)
	}

	if 3 == len(args) {
		onLoadEnd := args[2]
		evalTarget := CreateConsCell(onLoadEnd,
			CreateConsCell(CreateNil(), CreateNil()))

		return Eval(ctx, evalTarget, env)
	}

	return CreateNil(), nil
}

func NewFileReadLine() *Sexpression {
	return CreateSubroutine("file-read-line", _subr_file_read_line_Apply)
}
