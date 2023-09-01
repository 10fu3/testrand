package eval

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
	"testrand/reader/globalEnv"
)

var PutReceiveQueueMethod func(evnId string, reqId string, onReceive SExpression)

func StartReceiveServer(ctx context.Context) (func(), func(evnId string, reqId string, onReceive SExpression)) {
	m := sync.Map{}
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

		result, err = Eval(ctx, createSExpressionOnReceive, globalEnv.Get(sExpressionEnv.envId).(Environment))

		if err != nil {
			fmt.Println(err)
		}
	})

	return func() {
			router.Run(":4040")
		}, func(evnId string, reqId string, onReceive SExpression) {
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
