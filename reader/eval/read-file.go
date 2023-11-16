package eval

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
)

type _file_read struct{}

func (_ *_file_read) TypeId() string {
	return "subroutine.read-line-file"
}

func (_ *_file_read) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_file_read) String() string {
	return "#<subr read-file>"
}

func (_ *_file_read) IsList() bool {
	return false
}

func (l *_file_read) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_file_read) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {

	args, err := ToArray(arguments)
	if err != nil {
		return nil, err
	}

	if 1 > len(args) {
		return nil, errors.New("need arguments size is 1")
	}

	rawPath := args[0]
	if rawPath.TypeId() != "string" {
		return nil, errors.New("need arguments type is string")
	}

	path := rawPath.(Str).GetValue()

	fp, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	fileReader := bufio.NewReaderSize(fp, 64*1024)
	allReadStr := ""
	/*          第2引数はバッファサイズ。大きくするとある程度高速化する*/
	for {
		line_byte, _, err := fileReader.ReadLine()
		allReadStr += string(line_byte)

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return NewString(allReadStr), nil
}

func NewFileRead() SExpression {
	return &_file_read_line{}
}
