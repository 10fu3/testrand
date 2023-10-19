package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"testrand/reader/eval"
)

func main() {
	fmt.Println("light client")
	gin.SetMode(gin.ReleaseMode)
	ctx := context.Background()
	completed, addMethod := eval.StartReceiveServer(ctx)
	eval.PutReceiveQueueMethod = addMethod
	go func() {
		completed()
	}()
	env, err := eval.NewGlobalEnvironment()

	if err != nil {
		panic(err)
	}

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
