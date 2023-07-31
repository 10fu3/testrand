package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"testrand/reader/eval"
)

func main() {
	fmt.Println("light client")
	gin.SetMode(gin.ReleaseMode)
	completed, addMethod := eval.StartReceiveServer()
	eval.SetupPutReceiveQueueMethod(addMethod)
	go func() {
		completed()
	}()
	env := eval.NewGlobalEnvironment()

	stdin := bufio.NewReader(os.Stdin)
	read := eval.New(stdin)
	for {
		result, err := read.Read()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		result, err = eval.Eval(result, env)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if result != nil {
			fmt.Println(result)
		}
	}
}