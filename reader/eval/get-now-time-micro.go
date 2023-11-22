package eval

import (
	"context"
	"time"
)

func _subr_get_now_time_nano_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	return CreateInt(time.Now().UnixNano()), nil
}

func NewGetNowTimeNano() *Sexpression {
	return CreateSubroutine("get-now-time-nano", _subr_get_now_time_nano_Apply)
}
