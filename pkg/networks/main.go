package networks

import (
	"log"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/model/net"
	"github.com/loic-roux-404/crypto-bots/internal/nets/erc20"
)

// Nets map
var Nets = helpers.Map{
	erc20.ErcNetName: erc20.NewEth,
}

// GetNetwork in map
func GetNetwork(name string) net.Network {
	impl, err := helpers.GetInMap(Nets, name)
	if err != nil {
		log.Fatal(err)
	}

	netFn := impl.(func() net.Network)

	return netFn()
}
