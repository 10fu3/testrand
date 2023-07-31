package main

import (
	"bufio"
	"fmt"
	"os"
	"testrand/reader/eval"
)

func main() {

	eval.StartMockServer()
	eval.SetupPutReceiveQueueMethod(eval.StartReceiveServer())
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
