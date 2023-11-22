package eval

import (
	"context"
	"errors"
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

func (_ *_file_read) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	if 1 > argsLength {
		return nil, errors.New("need arguments size is 1")
	}

	rawPath := args[0]
	if rawPath.TypeId() != "string" {
		return nil, errors.New("need arguments type is string")
	}

	path := rawPath.(Str).GetValue()

	fileReader, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return NewString(string(fileReader)), nil
}

func NewFileRead() SExpression {
	return &_file_read{}
}
