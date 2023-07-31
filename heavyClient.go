package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testrand/reader/eval"
)

func main() {
	fmt.Println("heavy client")
	gin.SetMode(gin.ReleaseMode)
	eval.StartMockServer()
}
