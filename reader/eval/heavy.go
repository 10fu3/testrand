package eval

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func SendSExpression(sendSexp SExpression, onComplete SExpression, env Environment, port int) {
	reqId := uuid.NewString()
	PutReceiveQueueMethod(env.GetId(), reqId, onComplete)
	fromAddr := fmt.Sprintf("localhost:%d", port)
	values, err := json.Marshal(TaskAddRequest{
		Body: sendSexp.String(),
		From: &fromAddr,
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

	res, err := client.Post(fmt.Sprintf("http://localhost:8080/add-task/%s", reqId), "application/json", bytes.NewBuffer(values))
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

	if 1 == len(args) {
		SendSExpression(args[0], nil, env, 4040)
	}
	if 2 == len(args) {
		SendSExpression(args[0], args[1], env, 4040)
	}
	return nil, err
}

func NewHeavy() SExpression {
	return &_heavy{}
}
