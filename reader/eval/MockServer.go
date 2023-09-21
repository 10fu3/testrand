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
	Body string  `json:"body"`
	From *string `json:"from"`
}

func StartMockServer(ctx context.Context) {
	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		c.JSON(200, struct{ Message string }{Message: "OK"})
	})
	engine.GET("/routine-count", func(c *gin.Context) {
		c.JSON(200, struct{ Count int }{Count: runtime.NumGoroutine()})
	})
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, struct{ Status string }{Status: "OK"})
	})
	engine.POST("/add-task/:id", func(c *gin.Context) {
		requestId := c.Param("id")
		var req TaskAddRequest
		err := c.ShouldBind(&req)
		if requestId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ng",
				"message": "not allowed empty id",
			})
			return
		}
		if req.From == nil {
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

			_, err = http.Post(fmt.Sprintf("%s/receive/%s", *req.From, requestId), "application/json", sendBodyBuff)

			if err != nil {
				fmt.Println(err)
			}
			for i := 0; i < 5; i++ {
				if err == nil {
					break
				}
				time.Sleep(time.Second * 3)
				_, err = http.Post(fmt.Sprintf("%s/receive/%s", *req.From, requestId), "application/json", sendBodyBuff)
			}
		}()
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"id":     requestId,
		})
	})
	engine.Run(":8080")
}
