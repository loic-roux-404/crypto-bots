package net

import (
	"errors"
	"math/big"
	"path/filepath"

	"github.com/loic-roux-404/crypto-bots/internal/config"
	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/model/token"
	"github.com/loic-roux-404/crypto-bots/internal/model/wallet"
	"github.com/spf13/viper"
)

// Config of etherum like blockchain
type Config struct {
	NetName     string               `mapstructure:"network"`
	ChainName   string               `mapstructure:"chain"`
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

// NetChainType viper cnf id
const (
	NetChainType   = "chain"
	NetName        = "network"
)

var (
	// ErrIpcNotConfigured no ipc
	ErrIpcNotConfigured = errors.New("No IPC url configured")
)

// NewNetConfig create erc like blockchain handler
func NewNetConfig(defaults helpers.Map) (*Config, error) {
	if len(viper.GetString(NetName)) <= 0 {
		viper.Set(NetName, defaults[viper.GetString(NetChainType)])
	}

	cnfLoc := filepath.Join(NetChainType, viper.GetString(NetChainType))

	var cnfLocations = map[string]string{
		cnfLoc: viper.GetString(NetName),
	}

	cnf := &Config{}
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
