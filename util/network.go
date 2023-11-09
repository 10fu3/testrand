package util

import "net"

var localIP *string

// https://qiita.com/KEINOS/items/60c3bdbf2b0a28d961bf
func GetLocalIP() (string, error) {

	if localIP != nil {
		return *localIP, nil
	}

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr).IP.String()

	localIP = &localAddr

	return localAddr, nil
}
