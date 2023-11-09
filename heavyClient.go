package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"testrand/reader/eval"
)

func main() {
	fmt.Println("heavy client")
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	ctx := context.TODO()
	eval.StartMockServer(ctx)
}
