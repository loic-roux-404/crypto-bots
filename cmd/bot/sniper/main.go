package main

import (
	"math/big"

	"github.com/loic-roux-404/crypto-bots/cmd/bot/template"
	"github.com/loic-roux-404/crypto-bots/pkg/networks"
	"github.com/spf13/cobra"
)

func init() {
	template.InitNetCmd(&cobra.Command{
		Use:   "sniper",
		Short: "Sniper transaction",
	})
}

func main() {
	template.ExecuteNetCmd()
	n := networks.GetNetwork("erc20")

	n.Send("0x0", big.NewFloat(0.02))
}
