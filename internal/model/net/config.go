package net

import (
	"errors"
	"math/big"
	"path/filepath"

	"github.com/loic-roux-404/crypto-bots/internal/config"
	"github.com/loic-roux-404/crypto-bots/internal/model/token"
	"github.com/loic-roux-404/crypto-bots/internal/model/wallet"
	"github.com/spf13/viper"
)

// ERCConfig of etherum like blockchain
type ERCConfig struct {
	ManualFee   bool                 `mapstructure:"manualFee"`
	GasLimit    uint64               `mapstructure:"gasLimit"`
	GasPrice    int64                `mapstructure:"gasPrice"`
	Pass        string               `mapstructure:"pass"`
	Keystore    string               `mapstructure:"keystore"`
	Ipc         string               `mapstructure:"ipc"`
	Ws          string               `mapstructure:"Ws"`
	ChainID     int64                `mapstructure:"chainid"`
	FromAccount string               `mapstructure:"fromAccount"`
	Wallets     []wallet.ImportedKey `mapstructure:"wallets"`
}

// NetCnfID viper cnf id
const NetCnfID = "network"

var (
	// ErrIpcNotConfigured no ipc
	ErrIpcNotConfigured = errors.New("No IPC url configured")
)

// NewERCConfig create erc like blockchain handler
func NewERCConfig(networkID string, defaultNode string) (*ERCConfig, error) {

	if len(viper.GetString(NetCnfID)) <= 0 {
		viper.Set(NetCnfID, defaultNode)
	}

	cnfLoc := filepath.Join(NetCnfID, networkID)

	var cnfLocations = map[string]string{
		cnfLoc: viper.GetString(NetCnfID),
	}

	cnf := &ERCConfig{}
	// Search config in files
	config.Get(cnf, cnfLocations)
	// TODO override configs with flags

	if cnf.Ipc == "" {
		return nil, ErrIpcNotConfigured
	}
	// TODO switch case on fee system
	cnf.GasPrice = token.GweiToWei(big.NewInt(cnf.GasPrice)).Int64()

	return cnf, nil
}
