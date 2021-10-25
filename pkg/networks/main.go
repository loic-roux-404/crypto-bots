package networks

import (
	"log"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/inetworks"
)

var nets = helpers.Map{
	"eth": inetworks.NewEth,
}

// GetNetwork in map
func GetNetwork(name string) (inetworks.Network, error) {
    net, err := helpers.GetInMap(nets, name); if err != nil {
        log.Fatal(err)
	}

	netFn := net.(func()(inetworks.Network, error))

    return netFn()
}
