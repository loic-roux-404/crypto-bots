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

	n.Send("0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d", big.NewFloat(0.02))
}
