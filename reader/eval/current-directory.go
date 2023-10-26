package eval

import (
	"context"
	"os"
)

type _current_directory struct{}

func (_ *_current_directory) Type() string {
	return "subroutine.current-directory"
}

func (_ *_current_directory) String() string {
	return "#<subr current-directory>"
}

func (_ *_current_directory) IsList() bool {
	return false
}

func (l *_current_directory) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_current_directory) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {
	p, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	return NewString(p), nil
}
func NewCurrentDirectory() SExpression {
	return &_current_directory{}
}
