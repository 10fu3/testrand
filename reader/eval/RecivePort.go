package eval

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
	"sync"
	"testrand/config"
	"testrand/reader/globalEnv"
	"testrand/util"
)

var PutReceiveQueueMethod func(evnId string, reqId string, onReceive SExpression)

func StartReceiveServer(globalNamespaceId string, ctx context.Context) (func(), func(evnId string, reqId string, onReceive SExpression)) {
	m := sync.Map{}
	router := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	localIp, err := util.GetLocalIP()

	if err != nil {
		panic(err)
	}

	conf := config.Get()

	err = LoadBalancingRegisterForClient(struct {
		host  string
		port  string
		envId string
	}{
		host:  localIp,
		port:  "4040",
		envId: globalNamespaceId,
	}, struct {
		host string
		port string
	}{
		host: conf.ProxyHost,
		port: conf.ProxyPort,
	})

	if err == nil {

		router.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status": "ok",
			})
		})
		router.Post("/receive/:id", func(c *fiber.Ctx) error {
			var req struct {
				Result string `json:"result"`
			}
			reqId := c.Params("id")
			err := c.BodyParser(&req)
			if err != nil {
				fmt.Println(err.Error())
				m.Delete(reqId)
				return err
			}
			sample := strings.NewReader(fmt.Sprintf("%s\n", req.Result))
			read := New(bufio.NewReader(sample))
			result, err := read.Read()
			storedSExpressionEnv, ok := m.Load(reqId)

			if !ok {
				return errors.New("not found request id in callback store")
			}
			sExpressionEnv := storedSExpressionEnv.(*struct {
				onReceive SExpression
				envId     string
			})
			if sExpressionEnv.onReceive == nil {
				c.Status(http.StatusOK)
				m.Delete(reqId)
				return nil
			}
			createSExpressionOnReceive :=
				NewConsCell(sExpressionEnv.onReceive,
					NewConsCell(result,
						NewConsCell(NewNil(), NewNil())))

			targetEnv := globalEnv.Get(sExpressionEnv.envId)

			result, err = Eval(ctx, createSExpressionOnReceive, targetEnv.(Environment))

			if err != nil {
				m.Delete(reqId)
				fmt.Println(err)
			}
			m.Delete(reqId)
			globalEnv.Delete(sExpressionEnv.envId)
			c.Status(http.StatusOK)
			return nil
		})
	} else {
		fmt.Println("load balancing system is occurred error")
		fmt.Println(err)
	}

	return func() {
			router.Listen(fmt.Sprintf(":%s", conf.SelfOnCompletePort))
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
