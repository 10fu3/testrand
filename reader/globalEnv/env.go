package globalEnv

import "sync"

var envs = struct {
	m map[string]interface{}
	sync.Mutex
}{
	m: make(map[string]interface{}),
}

func Get(id string) interface{} {
	envs.Lock()
	defer envs.Unlock()
	return envs.m[id]
}

func Put(id string, env interface{}) {
	envs.Lock()
	defer envs.Unlock()
	envs.m[id] = env
}

func Delete(id string) {
	envs.Lock()
	defer envs.Unlock()
	delete(envs.m, id)
}
