package networking

import (
	"log"
	"net"
	"time"
)

const (
	wsProtocol  = "ws"
	sepProtocol = "://"
)

// Ping an url
func Ping(host string, port string) bool {
	timeout := time.Duration(1 * time.Second)

	_, err := net.DialTimeout("tcp", host+":"+port, timeout)

	if err != nil {
		log.Printf("%s %s %s\n", host, "not responding", err)
		return false
	}

	log.Printf("%s %s %s\n", host, "responding on port:", port)
	return true
}

// IsDomain check
func IsDomain(addr string) bool {
	print("addr from is dom", addr)
	ecounteredDot := false
	domain := addr[len(sepProtocol):]

	if domain[0] == '.' {
		return false
	}
	// Infinite subdomain validation
	for i, r := range domain {
		// Can't use dom..url
		if ecounteredDot && r == '.' {
			return false
		}

		if r != '.' {
			ecounteredDot = false
			continue
		}

		if i > 0 && r == '.' {
			ecounteredDot = true
		}
	}

	return true
}

func protocolURLOk(pr, addr string) bool {
	prLen := len(pr)
	if len(addr) <= prLen || addr[0:prLen] != pr {
		return false
	}

	sslLen := 0
	prPos := prLen - 1

	// SSL case
	if len(addr) >= prPos+1 && addr[prPos+1] == 's' {
		sslLen = 1
	}

	if len(addr) <= prLen+sslLen+len(sepProtocol) &&
		addr[prPos+sslLen:len(sepProtocol)] != sepProtocol {
		print(addr)
		print("aled", addr[prPos+sslLen:])
		return false
	}

	print("\naaaaaaaaaaa: ", addr[:prLen+sslLen])

	return IsDomain(addr[sslLen+len(sepProtocol):])
}

// IsWs websocket url detection
func IsWs(addr string) bool {
	return protocolURLOk(wsProtocol, addr)
}
