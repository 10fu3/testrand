package eval

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net"
	"net/http"
	"runtime"
	"strings"
	"testrand/config"
	"testrand/reader/globalEnv"
	"testrand/util"
	"time"
)

type TaskAddRequest struct {
	Body              *string `json:"body"`
	From              *string `json:"from"`
	GlobalNamespaceId *string `json:"global_namespace_id"`
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

	conf := config.Get()

	ramdomListener, _close := createListener()
	randomPort := fmt.Sprintf("%d", ramdomListener.Addr().(*net.TCPAddr).Port)
	_close()

	localIp, err := util.GetLocalIP()

	if err != nil {
		panic(err)
	}

	LoadBalancingRegisterForHeavy(struct {
		host string
		port string
	}{host: localIp, port: randomPort}, struct {
		host string
		port string
	}{host: conf.ProxyHost, port: conf.ProxyPort})

	engine := fiber.New()

	engine.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(struct {
			Message string `json:"message"`
		}{Message: "OK"})
	})
	engine.Get("/routine-count", func(c *fiber.Ctx) error {
		fmt.Printf("health check: %d\n", runtime.NumGoroutine())
		return c.JSON(struct {
			Count int `json:"count"`
		}{Count: runtime.NumGoroutine()})
	})
	engine.Get("/health", func(c *fiber.Ctx) error {
		fmt.Println("health check")
		return c.JSON(struct {
			Status string `json:"status"`
		}{Status: "OK"})
	})
	engine.Post("/add-task/:id", func(c *fiber.Ctx) error {
		requestId := c.Params("id")
		var req TaskAddRequest
		err := c.BodyParser(&req)
		if requestId == "" {
			return c.JSON(fiber.Map{
				"status":  "ng",
				"message": "not allowed empty id",
			})
		}
		if req.From == nil {
			return c.JSON(fiber.Map{
				"status":  "ng",
				"message": "not allowed empty port",
			})
		}
		if req.Body == nil {
			return c.JSON(fiber.Map{
				"status":  "ng",
				"message": "not allowed empty body",
			})
		}
		if req.GlobalNamespaceId == nil {
			return c.JSON(fiber.Map{
				"status":  "ng",
				"message": "not allowed empty session_id",
			})
		}
		go func() {
			if err != nil {
				fmt.Println("req err: " + err.Error())
				return
			}
			env, err := NewGlobalEnvironmentById(*req.GlobalNamespaceId)

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
			globalEnv.Delete(env.GetId())
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

			sendAddr := fmt.Sprintf("http://%s/receive/%s", *req.From, requestId)

			fmt.Println("sendAddr:", sendAddr)

			_, err = http.Post(sendAddr, "application/json", sendBodyBuff)

			if err != nil {
				fmt.Println(err)
			}
			for i := 0; i < 5; i++ {
				if err == nil {
					break
				}
				time.Sleep(time.Second * 3)
				_, err = http.Post(sendAddr, "application/json", sendBodyBuff)
			}

		}()
		return c.JSON(fiber.Map{
			"status": "ok",
			"id":     requestId,
		})
	})
	if err := engine.Listen(fmt.Sprintf(":%s", randomPort)); err != nil {
		panic(err)
	}
}
