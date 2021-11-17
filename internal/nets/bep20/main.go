package bep20

import (
	"github.com/loic-roux-404/crypto-bots/internal/model/net"
	"github.com/loic-roux-404/crypto-bots/internal/nets/erc20"
)

const (
	// BepNetName identifier
	BepNetName  = "bep20"
	// DefaultNet default net name (config file to load)
	DefaultNet = "testnet"
)

// NewBep binance smart chain module
func NewBep(cnf *net.Config) net.Network {
	return erc20.NewEth(cnf)
}
