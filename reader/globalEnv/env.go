package globalEnv

var envs = map[string]interface{}{}

func Get(id string) interface{} {
	return envs[id]
}

func Put(id string, env interface{}) {
	envs[id] = env
}
