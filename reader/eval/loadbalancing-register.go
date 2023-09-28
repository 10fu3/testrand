package eval

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func LoadBalancingRegister(host string, port int) {
	jsonContent := map[string]string{
		"machine_type": "reader",
		"from":         fmt.Sprintf("http://%s:%d", host, port),
	}
	jsonByte, err := json.Marshal(jsonContent)
	if err != nil {
		panic(err)
	}
	sendBodyBuff := bytes.NewBuffer(jsonByte)
	post, err := http.Post("http://localhost/register", "application/json", sendBodyBuff)
	if err != nil {
		panic(err)
	}
	fmt.Printf("regist result: %d\n", post.StatusCode)
}
