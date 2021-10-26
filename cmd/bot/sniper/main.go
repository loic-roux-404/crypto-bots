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
	sniperCmd.PersistentFlags().StringVar(
		&cfgFile,
		CnfFileID,
		CnfFileIDDefault,
		fmt.Sprintf("user config file (default is %s)", CnfFileIDDefault),
	)
	sniperCmd.PersistentFlags().StringP(
		net.NetCnfID,
		"n",
		NetCnfIDDefault,
		"ERC_20 like network id to load, default depending of chain type",
	)
	viper.BindPFlag(net.NetCnfID, sniperCmd.PersistentFlags().Lookup(net.NetCnfID))
}

func main() {
	sniperCmd.Execute()
	n, err := networks.GetNetwork("erc20")

	if err != nil {
		log.Fatalf("%s", err)
	}

	n.Send("0x36A130e8BD0fa0a39B92CfEEeCC8356EdbdD109e", token.NewBtcPair("ETH"), big.NewInt(1))
}
