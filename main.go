package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"testrand/reader/eval"
)

func main() {
	alreadyFlagChan := make(chan struct{})
	ctx := context.Background()
	env, err := eval.NewGlobalEnvironment()
	if err != nil {
		panic(err)
	}
	go func() {
		go func() {
			eval.StartMockServer(ctx)
		}()
		completed, addMethod := eval.StartReceiveServer(env.GetParentId(), ctx)
		eval.PutReceiveQueueMethod = addMethod
		alreadyFlagChan <- struct{}{}
		completed()
	}()
	<-alreadyFlagChan

	stdin := bufio.NewReader(os.Stdin)
	read := eval.New(stdin)
	for {
		result, err := read.Read()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		result, err = eval.Eval(ctx, result, env)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if result != nil {
			fmt.Println(result)
		}
	}
}
