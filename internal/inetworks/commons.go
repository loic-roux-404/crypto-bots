package inetworks

import (
	"math/big"

    "github.com/ethereum/go-ethereum/common"
)

// Network type
// Send to wallet and call smart contract
type Network interface {
	Send(address string, amount *big.Float) (hash common.Hash, err error)
	//Approve(address string) (hash common.Hash, err error)
	//Call(address string) (hash common.Hash, err error)

	Cancel(nonce *big.Int) (common.Hash, error)

	Update(nonce *big.Int, address string, amount *big.Float) (hash common.Hash, err error)
}
