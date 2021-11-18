package core

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/loic-roux-404/crypto-bots/internal/model/net"
)

const (
	// NetNameDefault default network name to load
	NetNameDefault = "ropsten"
	// NetChainDefault default blockchain
	NetChainDefault = "erc20"
	// CnfFileID config location
	CnfFileID = "config"
	// CnfFileIDDefault file
	CnfFileIDDefault = "config.yaml"
	// Flags
	keystoreID = "keystore"
	manualID   = "manualFee"
	chainid    = "chainid"
	pass       = "pass"
)

// genCmd generate command info
var genCmd *cobra.Command

// InitChainCmd use in init() function in your cli impl
// Dedicated to on chain system custom strategy
func InitChainCmd(infos *cobra.Command) {
	var (
		// CfgFile location
		CfgFile string
		// genCmd command to configure network and wallet
	)
	genCmd = infos
	// User Configuration
	genCmd.PersistentFlags().StringVar(
		&CfgFile,
		CnfFileID,
		CnfFileIDDefault,
		fmt.Sprintf("user config file (default is %s)", CnfFileIDDefault),
	)
	// Network type choice
	genCmd.PersistentFlags().StringP(
		net.NetChainType,
		"c",
		NetChainDefault,
		"Blockchain name, available : erc20 / bep20",
	)
	viper.BindPFlag(net.NetChainType, genCmd.PersistentFlags().Lookup(net.NetChainType))
	// Network environment choice
	genCmd.PersistentFlags().StringP(
		net.NetName,
		"n",
		"",
		"Network to load, default depending of chain type",
	)
	viper.BindPFlag(net.NetName, genCmd.PersistentFlags().Lookup(net.NetName))
	// Keystorefile location
	genCmd.PersistentFlags().StringP(
		keystoreID,
		"m",
		"",
		"Keystore location, json format file storing encrypted keys. create it from metamask",
	)
	genCmd.MarkPersistentFlagRequired(keystoreID)
	viper.BindPFlag(keystoreID, genCmd.PersistentFlags().Lookup(keystoreID))
	// Auto gas
	genCmd.PersistentFlags().Bool(manualID, false, "Disable automatic gas estimation")
	viper.BindPFlag(manualID, genCmd.PersistentFlags().Lookup(manualID))
	// Account password
	genCmd.PersistentFlags().StringP(
		pass,
		"p",
		"",
		"Account Password",
	)
	genCmd.MarkPersistentFlagRequired(pass)
	viper.BindPFlag(pass, genCmd.PersistentFlags().Lookup(pass))
}

// ExecuteNetCmd command for networks
func ExecuteNetCmd() error {
	return genCmd.Execute()
}
