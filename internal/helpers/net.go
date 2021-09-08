package helpers

import (
	"log"
	"net"
	"time"
)

func Ping(host string, port string) bool {
	timeout := time.Duration(1 * time.Second)

	_, err := net.DialTimeout("tcp", host + ":" + port, timeout)

	if err != nil {
		log.Printf("%s %s %s\n", host, "not responding", err.Error())
		return false
	}

	log.Printf("%s %s %s\n", host, "responding on port:", port)
	return true
}
