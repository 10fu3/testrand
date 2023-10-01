package eval

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type TaskAddRequest struct {
	Body      *string `json:"body"`
	From      *string `json:"from"`
	SessionId *string `json:"session_id"`
}

func createListener() (l net.Listener, close func()) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	return l, func() {
		_ = l.Close()
	}
}

func StartMockServer(ctx context.Context) {

	ramdomListener, _close := createListener()
	randomPort := ramdomListener.Addr().(*net.TCPAddr).Port
	_close()

	LoadBalancingRegister("localhost", randomPort)

	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		c.JSON(200, struct {
			Message string `json:"message"`
		}{Message: "OK"})
	})
	engine.GET("/routine-count", func(c *gin.Context) {
		fmt.Printf("health check: %d\n", runtime.NumGoroutine())
		c.JSON(200, struct {
			Count int `json:"count"`
		}{Count: runtime.NumGoroutine()})
	})
	engine.GET("/health", func(c *gin.Context) {
		fmt.Println("health check")
		c.JSON(200, struct {
			Status string `json:"status"`
		}{Status: "OK"})
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
		if req.Body == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ng",
				"message": "not allowed empty body",
			})
			return
		}
		if req.SessionId == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ng",
				"message": "not allowed empty session_id",
			})
			return
		}
		go func() {
			if err != nil {
				fmt.Println("req err: " + err.Error())
				return
			}
			env, err := NewGlobalEnvironmentById(*req.SessionId)

			if err != nil {
				panic(err)
			}

			input := strings.NewReader(fmt.Sprintf("%s\n", *req.Body))
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

			_, err = http.Post(fmt.Sprintf("http://%s/receive/%s", *req.From, requestId), "application/json", sendBodyBuff)

			if err != nil {
				fmt.Println(err)
			}
			for i := 0; i < 5; i++ {
				if err == nil {
					break
				}
				time.Sleep(time.Second * 3)
				_, err = http.Post(fmt.Sprintf("http://%s/receive/%s", *req.From, requestId), "application/json", sendBodyBuff)
			}
		}()
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"id":     requestId,
		})
	})
	engine.Run(fmt.Sprintf(":%d", randomPort))
}
