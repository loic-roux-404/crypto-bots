package net

import (
	"errors"
	"path/filepath"

	"github.com/loic-roux-404/crypto-bots/internal/services"
	"github.com/spf13/viper"
)

// ERCConfig of etherum like blockchain
type ERCConfig struct {
	ManualFee bool  `mapstructure:"manualFee"`
	GasLimit int64  `mapstructure:"gasLimit"`
	GasPrice int64  `mapstructure:"gasPrice"`
	Pass     string `mapstructure:"pass"`
	Keystore string `mapstructure:"keystore"`
	Ipc 	 string `mapstructure:"ipc"`
	ChainID  int64  `mapstructure:"chainid"`

}

// NetCnfID viper cnf id
const NetCnfID = "network"

var (
	// ErrIpcNotConfigured no ipc
	ErrIpcNotConfigured = errors.New("No IPC url configured")
)

// NewERCConfig create erc like blockchain handler
func NewERCConfig(networkID string, defaultNode string) (*ERCConfig, error)  {

	if len(viper.GetString(NetCnfID)) <= 0 {
		viper.Set(NetCnfID, defaultNode)
	}

	cnfLoc := filepath.Join(NetCnfID, networkID)

	var cnfLocations = map[string]string{
		cnfLoc: viper.GetString(NetCnfID),
	}

	cnf := &ERCConfig{}
	// Search config in files
	services.GetCnf(cnf, cnfLocations)
	// TODO override configs with flags

	if cnf.Ipc == "" {
		return nil, ErrIpcNotConfigured
	}

	return cnf, nil
}
