package networks

import (
    "github.com/loic-roux-404/crypto-bots/internal/helpers"
    "github.com/loic-roux-404/crypto-bots/internal/inetworks"
)

// NetworkMap network map
type NetworkMap map[string](func () (inetworks.Network, error))

var nets = NetworkMap{
	"eth": inetworks.NewEth,
}

// Get in map
func Get(name string) (inetworks.Network, error) {
    net, err := helpers.GetInMap(nets, name)

    return net.(inetworks.Network), err
}
