package main

import (
	"bufio"
	"fmt"
	"os"
	"testrand/reader"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	read := reader.New(stdin)

	for {
		result, err := read.Read()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println(result.String())
	}
}
