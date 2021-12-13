package strategy

import (
	"log"
	"math/big"

	"github.com/loic-roux-404/crypto-bots/pkg/networks"
)

func Send(toAddress, amount string) {
	n := networks.GetNetwork()

	f, ok := new(big.Float).SetString(amount)

	if !ok {
		log.Fatalf("Amount isn't a correct decimal number")
	}

	// "0xBe20D507fbdD6dAFd7a2ddE95c2d3f4618547F17"
	n.Send(toAddress, f)
}
