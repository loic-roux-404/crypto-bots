package brokers

import (
	"log"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/ibrokers"
)

var bs = helpers.FnMap {
    "binance": ibrokers.NewBinance(),
}

// Get in map
func Get(name string) (ibrokers.Broker, error) {
    broker, err := helpers.GetInMap(bs, name)

    if err != nil {
        log.Fatal(err)
    }

    return broker.(ibrokers.Broker), err
}
