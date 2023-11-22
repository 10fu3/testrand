package eval

import (
	"context"
	"runtime"
)

func _subr__force_gc_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	runtime.GC()
	return CreateBool(true), nil
}

func NewForceGC() *Sexpression {
	return CreateSubroutine("force-gc", _subr__force_gc_Apply)
}
