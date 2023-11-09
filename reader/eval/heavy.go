package eval

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testrand/config"
	"testrand/util"
)

func SendSExpression(sendSexp SExpression, onComplete SExpression, env Environment, fromHost string, fromPort string) {

	conf := config.Get()

	reqId := uuid.NewString()
	PutReceiveQueueMethod(env.GetId(), reqId, onComplete)
	fromAddr := fmt.Sprintf("%s:%s", fromHost, fromPort)
	sexpBody := sendSexp.String()
	id := env.GetParentId()
	values, err := json.Marshal(TaskAddRequest{
		Body:              &sexpBody,
		From:              &fromAddr,
		GlobalNamespaceId: &id,
	})

	transport := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{}
			return dialer.DialContext(ctx, "tcp4", addr)
		},
	}

	client := http.Client{
		Transport: &transport,
	}

	sendReqBody := map[string]string{
		"body": sendSexp.String(),
		"from": fromAddr,
	}
	sendReqBodyByte, err := json.Marshal(sendReqBody)
	send, err := http.Post(fmt.Sprintf("http://%s:%s/send-request", conf.ProxyHost, conf.ProxyPort), "application/json", bytes.NewBuffer(sendReqBodyByte))
	sendTargetResult := struct {
		Addr string `json:"addr"`
	}{}
	sendTargetResultByte, err := ioutil.ReadAll(send.Body)
	if err := json.Unmarshal(sendTargetResultByte, &sendTargetResult); err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Post(fmt.Sprintf("%s/add-task/%s", sendTargetResult.Addr, reqId), "application/json", bytes.NewBuffer(values))
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
}

type _heavy struct{}

func (_ *_heavy) Type() string {
	return "special_form.heavy"
}

func (_ *_heavy) String() string {
	return "#<syntax heavy>"
}

func (_ *_heavy) IsList() bool {
	return false
}

func (h *_heavy) Equals(sexp SExpression) bool {
	return h.Type() == sexp.Type()
}

func (_ *_heavy) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {
	args, err := ToArray(arguments)

	if err != nil {
		return nil, err
	}

	if ctx.Value("transaction") != nil {
		return nil, errors.New("transaction can not use in heavy")
	}

	conf := config.Get()
	ip, err := util.GetLocalIP()

	if err != nil {
		return nil, err
	}

	if 1 == len(args) {
		SendSExpression(args[0], nil, env, ip, conf.SelfOnCompletePort)
	}
	if 2 == len(args) {
		SendSExpression(args[0], args[1], env, ip, conf.SelfOnCompletePort)
	}
	return nil, err
}

func NewHeavy() SExpression {
	return &_heavy{}
}
