package net

import (
	"errors"
	"math/big"
	"path/filepath"

	"github.com/loic-roux-404/crypto-bots/internal/config"
	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/kecacc/fees"
	"github.com/spf13/viper"
)

// ImportedKey pair
type ImportedKey struct {
	Priv string `mapstructure:"priv"`
	Pass string `mapstructure:"pass"`
}

// Config of etherum like blockchain
type Config struct {
	NetName     string        `mapstructure:"network"`
	ChainName   string        `mapstructure:"chain"`
	ManualFee   bool          `mapstructure:"manualFee"`
	GasLimit    uint64        `mapstructure:"gasLimit"`
	GasPrice    int64         `mapstructure:"gasPrice"`
	Pass        string        `mapstructure:"pass"`
	Keystore    string        `mapstructure:"keystore"`
	Ipc         string        `mapstructure:"ipc"`
	Ws          string        `mapstructure:"ws"`
	ChainID     int64         `mapstructure:"chainid"`
	FromAddress string        `mapstructure:"fromAddress"`
	Wallets     []ImportedKey `mapstructure:"wallets"`
}

const (
	// NetChainType viper cnf id
	NetChainType = "chain"
	// NetName cnf id
	NetName = "network"
)

var (
	// ErrIpcNotConfigured no ipc
	ErrIpcNotConfigured = errors.New("no IPC url configured")
)

// NewNetConfig create erc like blockchain handler
func NewNetConfig(defaults helpers.SimpleMap) (*Config, error) {
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
	cnf.GasPrice = fees.GweiToWei(big.NewInt(cnf.GasPrice)).Int64()

	return cnf, nil
}
