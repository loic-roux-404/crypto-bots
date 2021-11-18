package main

import (
	"math/big"

	"github.com/loic-roux-404/crypto-bots/cmd/bot/core"
	"github.com/loic-roux-404/crypto-bots/pkg/networks"
	"github.com/spf13/cobra"
)

func init() {
	core.InitChainCmd(&cobra.Command{
		Use:   "sniper",
		Short: "Sniper transaction",
	})
}

func main() {
	core.ExecuteNetCmd()
	n := networks.GetNetwork()

	n.Send("0xBe20D507fbdD6dAFd7a2ddE95c2d3f4618547F17", big.NewFloat(0.02))
}
