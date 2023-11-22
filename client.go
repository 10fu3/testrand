package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"testrand/reader/eval"
)

func main() {
	fmt.Println("light client")
	ctx := context.Background()

	env, err := eval.NewGlobalEnvironment()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("cant use etcd and super global environment")
	}

	completed, addMethod := eval.StartReceiveServer(env.GetEnvParentId(), ctx)
	eval.PutReceiveQueueMethod = addMethod
	go func() {
		completed()
	}()

	stdin := bufio.NewReader(os.Stdin)
	read := eval.New(stdin)
	for {
		result, readErr := read.Read()
		if readErr != nil {
			fmt.Println(readErr.Error())
			continue
		}
		result, runErr := eval.Eval(ctx, result, env)
		if runErr != nil {
			fmt.Println(runErr)
			continue
		}
		fmt.Println(result)
	}
}
