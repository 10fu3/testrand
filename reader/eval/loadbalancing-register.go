package eval

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func LoadBalancingRegisterForHeavy(self struct {
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
		"machine_type": "heavy",
		"from":         fmt.Sprintf("http://%s:%s", self.host, self.port),
	}
	jsonByte, err := json.Marshal(jsonContent)
	if err != nil {
		panic(err)
	}
	sendBodyBuff := bytes.NewBuffer(jsonByte)
	post, err := http.Post(fmt.Sprintf("http://%s:%s/register-heavy", loadBalancer.host, loadBalancer.port), "application/json", sendBodyBuff)
	if err != nil {
		panic(err)
	}
	fmt.Printf("regist result: %d\n", post.StatusCode)
}

func LoadBalancingRegisterForClient(self struct {
	host  string
	port  string
	envId string
}, loadBalancer struct {
	host string
	port string
}) {

	if self.host == loadBalancer.host {
		loadBalancer.host = "localhost"
	}

	jsonContent := map[string]string{
		"machine_type":        "client",
		"from":                fmt.Sprintf("http://%s:%s", self.host, self.port),
		"global_namespace_id": self.envId,
	}
	jsonByte, err := json.Marshal(jsonContent)
	if err != nil {
		panic(err)
	}
	sendBodyBuff := bytes.NewBuffer(jsonByte)
	post, err := http.Post(fmt.Sprintf("http://%s:%s/register-client", loadBalancer.host, loadBalancer.port), "application/json", sendBodyBuff)
	if err != nil {
		panic(err)
	}
	fmt.Printf("regist result: %d\n", post.StatusCode)
}
