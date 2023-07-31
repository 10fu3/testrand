package reader

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
	"testrand/reader/eval"
)

var PutReceiveQueueMethod = func(evnId string, reqId string, onReceive eval.SExpression) {
	return
}

func SetupPutReceiveQueueMethod(f func(evnId string, reqId string, onReceive eval.SExpression)) {
	PutReceiveQueueMethod = f
}

func StartReceiveServer() func(evnId string, reqId string, onReceive eval.SExpression) {
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
			read := New(bufio.NewReader(sample))
			result, err := read.Read()
			fmt.Println(result)

		})
	}()

	return func(evnId string, reqId string, onReceive eval.SExpression) {
		m.Store(reqId, struct {
			onReceive eval.SExpression
			envId     string
		}{
			onReceive: onReceive,
			envId:     evnId,
		})
	}
}
