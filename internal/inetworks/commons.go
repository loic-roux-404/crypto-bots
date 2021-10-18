package inetworks

import (
	"math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/loic-roux-404/crypto-bots/internal/model/token"
)

// Network type
// Send to wallet and call smart contract
type Network interface {
	Send(address string, pair token.Pair, amount *big.Int) (hash common.Hash, err error)
	Call(address string) (hash common.Hash, err error)
}
