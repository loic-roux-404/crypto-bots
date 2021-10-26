package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/loic-roux-404/crypto-bots/pkg/networks"
	"github.com/loic-roux-404/crypto-bots/internal/model/token"
	"github.com/loic-roux-404/crypto-bots/internal/model/net"
)

const (
	// NetCnfIDDefault default network name to load
	NetCnfIDDefault = "ropsten"
	// CnfFileID config location
	CnfFileID = "config"
	// CnfFileIDDefault file
	CnfFileIDDefault = "config.yaml"

	memonicID = "memonic"
	manualID = "manualFee"
)

var (
	// Used for flags.
	// network string
	cfgFile string

	sniperCmd = &cobra.Command{
		Use:   "sniper",
		Short: "Sniping wallet strategy bot",
		Long: `Use this bot to make simple transaction / contract calls
on a ERC-20 like blockchain or configure a more sophisticated strategy`,
	}
)

func init() {
	// User Configuration
	sniperCmd.PersistentFlags().StringVar(
		&cfgFile,
		CnfFileID,
		CnfFileIDDefault,
		fmt.Sprintf("user config file (default is %s)", CnfFileIDDefault),
	)
	// Network choice
	sniperCmd.PersistentFlags().StringP(
		net.NetCnfID,
		"n",
		NetCnfIDDefault,
		"ERC_20 like network id to load, default depending of chain type",
	)
	viper.BindPFlag(net.NetCnfID, sniperCmd.PersistentFlags().Lookup(net.NetCnfID))
	// Memonic
	sniperCmd.PersistentFlags().StringP(
		memonicID,
		"m",
		"",
		"Seed phrase (12 words)",
	)
	viper.BindPFlag(memonicID, sniperCmd.PersistentFlags().Lookup(memonicID))
	// Auto gas
	sniperCmd.PersistentFlags().Bool(manualID, false, "Disable automatic gas estimation")
	viper.BindPFlag(manualID, sniperCmd.PersistentFlags().Lookup(manualID))
}

func main() {
	sniperCmd.Execute()
	n, err := networks.GetNetwork("erc20")

	if err != nil {
		log.Fatalf("%s", err)
	}

	h, err := n.Send("0x36A130e8BD0fa0a39B92CfEEeCC8356EdbdD109e", token.NewBtcPair("ETH"), big.NewInt(1))

	if err != nil {
		log.Printf("Error: %s", err)
	} else {
		log.Printf("Sucessfuly sent tx: %s", h)
	}
}
