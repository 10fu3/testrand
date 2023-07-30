package eval

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
	"testrand/reader"
)

func StartReceiveServer() func(reqId string, onReceive SExpression) {
	m := sync.Map{}
	go func() {
		router := gin.Default()
		router.POST("/receive/:id", func(c *gin.Context) {
			var req struct {
				Result string `json:"result"`
			}
			_ = c.Param("id")
			err := c.BindJSON(&req)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			sample := strings.NewReader(fmt.Sprintf("%s\n", req.Result))
			read := reader.New(bufio.NewReader(sample))
			result, err := read.Read()
			fmt.Println(result)
		})
	}()

	return func(reqId string, onReceive SExpression) {
		m.Store(reqId, onReceive)
	}
}
