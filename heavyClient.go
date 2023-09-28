package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"testrand/reader/eval"
)

func main() {
	fmt.Println("heavy client")
	gin.SetMode(gin.ReleaseMode)
	ctx := context.TODO()
	eval.StartMockServer(ctx)
}
