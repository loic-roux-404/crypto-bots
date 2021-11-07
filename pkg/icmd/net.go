package icmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	manualID   = "manualFee"
	chainid    = "chainid"
	pass       = "pass"
)

// genCmd generate command info
var genCmd *cobra.Command

// InitNetCmd use in init() function in your cli impl
func InitNetCmd(infos *cobra.Command) {
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
	// Network choice
	genCmd.PersistentFlags().StringP(
		net.NetCnfID,
		"n",
		NetCnfIDDefault,
		"ERC_20 like network id to load, default depending of chain type",
	)
	viper.BindPFlag(net.NetCnfID, genCmd.PersistentFlags().Lookup(net.NetCnfID))
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
	// Chainid
	genCmd.PersistentFlags().Int16P(
		chainid,
		"i",
		3,
		"Chain id",
	)
	genCmd.MarkPersistentFlagRequired(chainid)
	viper.BindPFlag(chainid, genCmd.PersistentFlags().Lookup(chainid))
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
