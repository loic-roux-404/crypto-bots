package main

import "github.com/loic-roux-404/crypto-bots/pkg/networks"

func main() {
	n, _ := networks.GetNetwork("eth")
	println(&n)
	//n.Send()
}
