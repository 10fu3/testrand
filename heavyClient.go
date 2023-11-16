package main

import (
	"context"
	"fmt"
	"testrand/reader/eval"
)

func main() {
	fmt.Println("heavy client")
	ctx := context.TODO()
	eval.StartMockServer(ctx)
}
