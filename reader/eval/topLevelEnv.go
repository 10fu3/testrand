package eval

import (
	"sync"
)

var envs = struct {
	m map[string]*Sexpression
	sync.RWMutex
}{
	m: make(map[string]*Sexpression),
}

func TopLevelEnvGet(id string) *Sexpression {
	envs.RLock()
	defer envs.RUnlock()
	return envs.m[id]
}

func TopLevelEnvPut(id string, env *Sexpression) {
	envs.Lock()
	defer envs.Unlock()
	envs.m[id] = env
}

func TopLevelEnvDelete(id string) {
	envs.Lock()
	delete(envs.m, id)
	envs.Unlock()
}
