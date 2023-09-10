package eval

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type TaskAddRequest struct {
	Body        string `json:"body"`
	ReceivePort *int   `json:"receive_port"`
}

func StartMockServer(ctx context.Context) {
	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		c.JSON(200, struct{ Message string }{Message: "OK"})
	})
	engine.GET("/routine-count", func(c *gin.Context) {
		c.JSON(200, struct{ Count int }{Count: runtime.NumGoroutine()})
	})
	engine.POST("/add-task/:id", func(c *gin.Context) {
		requestId := c.Param("id")
		from := c.RemoteIP()
		var req TaskAddRequest
		err := c.ShouldBind(&req)
		if requestId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ng",
				"message": "not allowed empty id",
			})
			return
		}
		if req.ReceivePort == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ng",
				"message": "not allowed empty port",
			})
			return
		}
		go func() {
			if err != nil {
				fmt.Println("req err: " + err.Error())
				return
			}
			env := NewGlobalEnvironment()
			input := strings.NewReader(fmt.Sprintf("%s\n", req.Body))
			read := New(bufio.NewReader(input))
			readSexp, err := read.Read()
			if err != nil {
				fmt.Println("read err: " + err.Error())
				return
			}
			result, err := Eval(ctx, readSexp, env)
			if err != nil {
				fmt.Println(err)
				return
			}
			sendBody := struct {
				Result string `json:"result"`
			}{
				Result: result.String(),
			}
			sendBodyBytes, err := json.Marshal(&sendBody)
			sendBodyBuff := bytes.NewBuffer(sendBodyBytes)

			_, err = http.Post(fmt.Sprintf("http://%s:%d/receive/%s", from, *req.ReceivePort, requestId), "application/json", sendBodyBuff)

			if err != nil {
				fmt.Println(err)
			}
			for i := 0; i < 5; i++ {
				if err == nil {
					break
				}
				time.Sleep(time.Second * 3)
				_, err = http.Post(fmt.Sprintf("%s:", from), "application/json", sendBodyBuff)
			}
		}()
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"id":     requestId,
		})
	})
	engine.Run(":8080")
}
