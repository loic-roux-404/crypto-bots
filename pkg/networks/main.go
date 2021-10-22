package networks

import (
	"log"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/inetworks"
)

// NetworkMap network map
type NetworkMap map[string](inetworks.Network)

var nets = helpers.FnMap{
	"eth": inetworks.NewEth(),
}

// GetNetwork in map
func GetNetwork(name string) (inetworks.Network, error) {
    net, err := helpers.GetInMap(nets, name)

    return net.(inetworks.Network), err
}
