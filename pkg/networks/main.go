package networks

import (
	"log"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/model/net"
	"github.com/loic-roux-404/crypto-bots/internal/nets/bep20"
	"github.com/loic-roux-404/crypto-bots/internal/nets/erc20"
)

var defaults = helpers.SimpleMap{
	erc20.ErcNetName: "ropsten",
	bep20.BepNetName: "testnet",
}

// Nets map
var Nets = helpers.SimpleMap{
	erc20.ErcNetName: erc20.NewEth,
	bep20.BepNetName: bep20.NewBep,
}

// GetNetwork in map
func GetNetwork() net.Network {
	cnf, err := net.NewNetConfig(defaults)
	if err != nil {
		log.Fatal(err)
	}

	impl, err := helpers.GetInMap(Nets, cnf.ChainName)
	if err != nil {
		log.Fatal(err)
	}

	netFn := impl.(func(*net.Config) net.Network)

	return netFn(cnf)
}
