package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/loic-roux-404/crypto-bots/pkg/networks"
	"github.com/loic-roux-404/crypto-bots/internal/model/net"
)

const (
	// NetCnfIDDefault default network name to load
	NetCnfIDDefault = "ropsten"
	// CnfFileID config location
	CnfFileID = "config"
	// CnfFileIDDefault file
	CnfFileIDDefault = "config.yaml"
	// Flags
	keystoreID = "keystore"
	manualID = "manualFee"
	chainid = "chainid"
	pass = "pass"
)

var (
	// CfgFile location
	CfgFile string

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
		&CfgFile,
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
	// Keystorefile location
	sniperCmd.PersistentFlags().StringP(
		keystoreID,
		"m",
		"",
		"Keystore location, json format file storing encrypted keys. create it from metamask",
	)
	sniperCmd.MarkPersistentFlagRequired(keystoreID)
	viper.BindPFlag(keystoreID, sniperCmd.PersistentFlags().Lookup(keystoreID))
	// Auto gas
	sniperCmd.PersistentFlags().Bool(manualID, false, "Disable automatic gas estimation")
	viper.BindPFlag(manualID, sniperCmd.PersistentFlags().Lookup(manualID))
	// Chainid
	sniperCmd.PersistentFlags().Int16P(
		chainid,
		"i",
		3,
		"Chain id",
	)
	sniperCmd.MarkPersistentFlagRequired(chainid)
	viper.BindPFlag(chainid, sniperCmd.PersistentFlags().Lookup(chainid))
	// Account password
	sniperCmd.PersistentFlags().StringP(
		pass,
		"p",
		"",
		"Account Password",
	)
	sniperCmd.MarkPersistentFlagRequired(pass)
	viper.BindPFlag(pass, sniperCmd.PersistentFlags().Lookup(pass))

}

func main() {
	sniperCmd.Execute()
	n, err := networks.GetNetwork("erc20")

	if err != nil {
		log.Fatalf("%s", err)
	}

	h, err := n.Send("0xE216378C0ed702D66e09D4aDBE9548C52604eB6E", big.NewFloat(0.02))

	if err != nil {
		log.Printf("Error: %s", err)
	} else {
		log.Printf("Sucessfuly sent tx: %s", h)
	}
}
