package inetworks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/loic-roux-404/crypto-bots/internal/model/store"
)

// Network type
// Send to wallet and call smart contract
// All method panic an err if one occurs
type Network interface {
	Send(address string, amount *big.Float) (hash common.Hash)

	Cancel(nonce *big.Int) (common.Hash)

	Update(nonce *big.Int, address string, amount *big.Float) (hash common.Hash)

	Load(address string, loadFn store.LoadFn) *store.Contract

	Deploy(input string, storeDeployFn store.DeployFn) *store.Contract
}
