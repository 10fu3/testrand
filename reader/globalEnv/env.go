package globalEnv

import "sync"

var envs = struct {
	m map[string]interface{}
	sync.RWMutex
}{
	m: make(map[string]interface{}),
}

func Get(id string) interface{} {
	envs.RLock()
	defer envs.RUnlock()
	return envs.m[id]
}

func Put(id string, env interface{}) {
	envs.Lock()
	defer envs.Unlock()
	envs.m[id] = env
}

func Delete(id string) {
	envs.Lock()
	delete(envs.m, id)
	envs.Unlock()
}
