package net

import (
	"errors"
	"path/filepath"

	"github.com/loic-roux-404/crypto-bots/internal/services"
	"github.com/spf13/viper"
)

// ERCConfig of etherum like blockchain
type ERCConfig struct {
	GasLimit int64 `mapstructure:"gasLimit"`
	GasPrice int64 `mapstructure:"gasPrice"`
	Memonic  string `mapstructure:"MEMONIC"`
	Keystore string `mapstructure:"keystore"`
	Ipc 	 string `mapstructure:"ipc"`
}

// NetCnfID viper cnf id
const NetCnfID = "network"

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
		return nil, errors.New("No IPC url configured")
	}

	return cnf, nil
}
