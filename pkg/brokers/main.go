package brokers

import (
	"github.com/loic-roux-404/crypto-bots/internal/ibrokers"
	"github.com/loic-roux-404/crypto-bots/internal/helpers"
)

// BrokerMap type
type BrokerMap map[string](func () (ibrokers.Broker, error))

var bs = BrokerMap {
    "binance": ibrokers.NewBinance,
}

// Get in map
func Get(name string) (ibrokers.Broker, error) {
    broker, err := helpers.GetInMap(bs, name)

    return broker.(ibrokers.Broker), err
}
