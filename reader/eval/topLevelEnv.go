package eval

import (
	"sync"
)

var envs = struct {
	m map[string]Environment
	sync.RWMutex
}{
	m: make(map[string]Environment),
}

func TopLevelEnvGet(id string) Environment {
	envs.RLock()
	defer envs.RUnlock()
	return envs.m[id]
}

func TopLevelEnvPut(id string, env Environment) {
	envs.Lock()
	defer envs.Unlock()
	envs.m[id] = env
}

func TopLevelEnvDelete(id string) {
	envs.Lock()
	delete(envs.m, id)
	envs.Unlock()
}
