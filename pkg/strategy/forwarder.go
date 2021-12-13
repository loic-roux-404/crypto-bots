package strategy

import (
	"log"

	"github.com/loic-roux-404/crypto-bots/internal/model/tx"
	"github.com/loic-roux-404/crypto-bots/pkg/networks"
)

func Forwarder(toAddress string) {
	n := networks.GetNetwork()

	// "0xBe20D507fbdD6dAFd7a2ddE95c2d3f4618547F17"

	n.SubscribeCurrent().RunEventLoop(func(tx *tx.Adapter) {
		log.Printf("Info: Forwarding fee to %s", toAddress)
		tx.Log()
	})
}
