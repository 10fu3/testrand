package main

import (
	"bufio"
	"fmt"
	"os"
	"testrand/reader"
	"testrand/reader/eval"
)

func main() {

	eval.StartMockServer()

	env := eval.NewGlobalEnvironment()

	stdin := bufio.NewReader(os.Stdin)
	read := reader.New(stdin)
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
