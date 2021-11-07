package net

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/loic-roux-404/crypto-bots/internal/model/store"
	"github.com/loic-roux-404/crypto-bots/internal/watcher"
)

// Network type
// Send to wallet and call smart contract
// All method panic an err if one occurs
type Network interface {
	Send(address string, amount *big.Float) common.Hash

	Cancel(nonce *big.Int) common.Hash

	Update(nonce *big.Int, address string, amount *big.Float) common.Hash

	Load(address string, loadFn store.LoadFn) interface{}

	Deploy(input string, storeDeployFn store.DeployFn) interface{}

	Subscribe(address string) watcher.WatcherSc

	SubscribeCurrent() watcher.WatcherAcc
}
