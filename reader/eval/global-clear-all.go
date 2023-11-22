package eval

import "context"

func _subr_global_clear_all_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	if err := env._env_superGlobalEnv.ClearAll(); err != nil {
		return CreateNil(), err
	}
	return CreateNil(), nil
}

func NewGlobalClearAll() *Sexpression {
	return CreateSubroutine("global-clear-all", _subr_global_clear_all_Apply)
}
