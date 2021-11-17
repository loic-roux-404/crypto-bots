package net

import (
	"math/big"

	"github.com/loic-roux-404/crypto-bots/internal/kecacc/store"
	"github.com/loic-roux-404/crypto-bots/internal/model/sub"
)

// Network type
// Send to wallet and call smart contract
// All method panic if error occurs
type Network interface {
	// Send a simple tx to address
	Send(address string, amount *big.Float) string
	// Cancel a pending tx
	Cancel(nonce *big.Int) string
	// Update a pending tx
	Update(nonce *big.Int, address string, amount *big.Float) string
	// Load smart contract with erc20 net
	Load(address string, loadFn store.LoadFn) interface{}
	// Deploy smart contract from abi generated function
	Deploy(input string, storeDeployFn store.DeployFn) interface{}
	// Smart contract events
	Subscribe(address string) sub.Sc
	// Account subscribe
	SubscribeCurrent() sub.Acc
}
