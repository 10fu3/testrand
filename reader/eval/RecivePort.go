package eval

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
	"testrand/reader/globalEnv"
)

var PutReceiveQueueMethod = func(evnId string, reqId string, onReceive SExpression) {
	return
}

func SetupPutReceiveQueueMethod(f func(evnId string, reqId string, onReceive SExpression)) {
	PutReceiveQueueMethod = f
}

func StartReceiveServer() func(evnId string, reqId string, onReceive SExpression) {
	m := sync.Map{}
	go func() {
		router := gin.Default()
		router.POST("/receive/:id", func(c *gin.Context) {
			var req struct {
				Result string `json:"result"`
			}
			reqId := c.Param("id")
			err := c.BindJSON(&req)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			sample := strings.NewReader(fmt.Sprintf("%s\n", req.Result))
			read := New(bufio.NewReader(sample))
			result, err := read.Read()
			fmt.Println(result)
			storedSExpressionEnv, ok := m.Load(reqId)

			if !ok {
				return
			}
			sExpressionEnv := storedSExpressionEnv.(*struct {
				onReceive SExpression
				envId     string
			})
			if sExpressionEnv.onReceive == nil {
				return
			}
			createSExpressionOnReceive :=
				NewConsCell(sExpressionEnv.onReceive,
					NewConsCell(result,
						NewConsCell(NewNil(), NewNil())))

			result, err = Eval(createSExpressionOnReceive, globalEnv.Get(sExpressionEnv.envId).(Environment))

			if err != nil {
				fmt.Println(err)
			}
		})
		router.Run(":4040")
	}()

	return func(evnId string, reqId string, onReceive SExpression) {
		stored := struct {
			onReceive SExpression
			envId     string
		}{
			onReceive: onReceive,
			envId:     evnId,
		}

		m.Store(reqId, &stored)
	}
}
