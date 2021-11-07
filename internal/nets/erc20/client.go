package erc20

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// NodeClients store
type NodeClients struct {
	Ws  *rpc.Client
	Ipc *rpc.Client
}

// NewClients for erc20
// TODO save memory with
// - ws non connect
// - store instances
func NewClients(ipc string, ws string) (nc *NodeClients, err error) {
	nc = &NodeClients{}
	nc.Ipc, err = rpc.Dial(ipc)
	if err != nil {
		return nil, err
	}

	nc.Ws, err = rpc.Dial(ws)
	if err != nil {
		return nil, err
	}

	return nc, nil
}

// GethRPC gethclient
func (nc *NodeClients) GethRPC() *gethclient.Client {
	return gethclient.New(nc.Ipc)
}

// EthRPC ethclient
func (nc *NodeClients) EthRPC() *ethclient.Client {
	return ethclient.NewClient(nc.Ipc)
}

// EthWs ethclient for websocket
func (nc *NodeClients) EthWs() *ethclient.Client {
	return ethclient.NewClient(nc.Ws)
}

// GethWs ethclient for websocket
func (nc *NodeClients) GethWs() *gethclient.Client {
	return gethclient.New(nc.Ws)
}
