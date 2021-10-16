package networks

import (
    "github.com/loic-roux-404/crypto-bots/internal/helpers"
    "github.com/loic-roux-404/crypto-bots/internal/inetworks"
)

// NetworkMap network map
type NetworkMap map[string](func () (inetworks.Network, error))

var nets = map[string]inetworks.Network{
	"eth": inetworks.NewEth,
}

// Get in map
func Get(name string) (inetworks.Network, error) {
    broker, err := helpers.GetInMap(nets, name)

    return broker.(inetworks.Network), err
}
