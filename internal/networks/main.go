package networks

import (
	"math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/loic-roux-404/pump-bot/internal/model/symbol"
)

// Network type
// Send to wallet and call smart contract
type Network interface {
	Send(address string, amount *big.Int, symbol symbol.Pair) (hash common.Hash, err error)
	Call(address string) (hash *common.Hash, err error)
}

var eth, _ = NewEth()

var brokers = map[string]Network{
	"eth": eth,
}

// GetBroker in map
func GetBroker(name string) Network {
    return brokers[name]
}
