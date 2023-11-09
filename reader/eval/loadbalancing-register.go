package eval

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func LoadBalancingRegister(self struct {
	host string
	port string
}, loadBalancer struct {
	host string
	port string
}) {

	if self.host == loadBalancer.host {
		loadBalancer.host = "localhost"
	}

	jsonContent := map[string]string{
		"machine_type": "reader",
		"from":         fmt.Sprintf("http://%s:%d", self.host, self.port),
	}
	jsonByte, err := json.Marshal(jsonContent)
	if err != nil {
		panic(err)
	}
	sendBodyBuff := bytes.NewBuffer(jsonByte)
	post, err := http.Post(fmt.Sprintf("http://%s:%d", loadBalancer.host, loadBalancer.port), "application/json", sendBodyBuff)
	if err != nil {
		panic(err)
	}
	fmt.Printf("regist result: %d\n", post.StatusCode)
}
