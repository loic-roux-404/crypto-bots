package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/loic-roux-404/crypto-bots/pkg/networks"
	"github.com/loic-roux-404/crypto-bots/internal/model/token"
)

const ( 
	// NetCnfID viper cnf id
	NetCnfID = "network"
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
		// Run: args.ShowFullHelp,
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
		NetCnfID, 
		"n", 
		"", 
		"ERC_20 like network id to load, default depending of chain type",
	)
	viper.BindPFlag(NetCnfID, sniperCmd.PersistentFlags().Lookup(NetCnfID))
	viper.SetDefault(NetCnfID, NetCnfIDDefault)
}

func main() {
	sniperCmd.Execute()
	n, err := networks.GetNetwork("eth")

	if err != nil {
		log.Fatalf("%s", err)
	}

	n.Send("0x36A130e8BD0fa0a39B92CfEEeCC8356EdbdD109e", token.NewBtcPair("ETH"), big.NewInt(1))
}
