package inetworks

import (
	"math/big"

    "github.com/ethereum/go-ethereum/common"
)

// Network type
// Send to wallet and call smart contract
type Network interface {
	// Send method
	// address is blockchain dependant cryptography
	// amount is calculated with the native blockchain token
	// account is a key derivation path
	Send(address string, amount *big.Float) (hash common.Hash, err error)
	//Approve(address string) (hash common.Hash, err error)
	//Call(address string) (hash common.Hash, err error)
}
