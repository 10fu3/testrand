package eval

import (
	"context"
	"os"
)

type _current_directory struct{}

func _sub_current_directory_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	p, err := os.Getwd()

	if err != nil {
		return CreateNil(), err
	}

	return CreateString(p), nil
}
func NewCurrentDirectory() *Sexpression {
	return CreateSubroutine("current-directory", _sub_current_directory_Apply)
}
