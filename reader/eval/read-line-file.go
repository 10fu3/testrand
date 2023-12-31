package eval

import (
	"bufio"
	"context"
	"errors"
	"os"
)

type _file_read_line struct{}

func (_ *_file_read_line) TypeId() string {
	return "subroutine.read-line-file"
}

func (_ *_file_read_line) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_file_read_line) String() string {
	return "#<subr read-line-file>"
}

func (_ *_file_read_line) IsList() bool {
	return false
}

func (l *_file_read_line) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_file_read_line) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	// 1st filepath string
	// 2nd on-load-line function
	// 3rd on-load-end function

	if 2 > argsLength {
		return nil, errors.New("need arguments size is 1")
	}

	rawPath := args[0]
	if rawPath.TypeId() != "string" {
		return nil, errors.New("need arguments type is string")
	}

	path := rawPath.(Str).GetValue()

	onLoadLine := args[1]

	var fp *os.File

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		evalTarget := NewConsCell(onLoadLine,
			NewConsCell(NewString(scanner.Text()),
				NewConsCell(NewNil(), NewNil())))

		Eval(ctx, evalTarget, env)
	}

	if 3 == len(args) {
		onLoadEnd := args[2]
		evalTarget := NewConsCell(onLoadEnd,
			NewConsCell(NewNil(), NewNil()))

		return Eval(ctx, evalTarget, env)
	}

	return NewConsCell(NewNil(), NewNil()), nil
}

func NewFileReadLine() SExpression {
	return &_file_read_line{}
}
