package clients

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// NodeErcClients store
type NodeErcClients struct {
	Ws  *rpc.Client
	Ipc *rpc.Client
}

// NewClients for erc20
// TODO save memory with
// - ws non connect
// - store instances
func NewClients(ipc string, ws string) (nc *NodeErcClients, err error) {
	nc = &NodeErcClients{}
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
func (nc *NodeErcClients) GethRPC() *gethclient.Client {
	return gethclient.New(nc.Ipc)
}

// EthRPC ethclient
func (nc *NodeErcClients) EthRPC() *ethclient.Client {
	return ethclient.NewClient(nc.Ipc)
}

// EthWs ethclient for websocket
func (nc *NodeErcClients) EthWs() *ethclient.Client {
	return ethclient.NewClient(nc.Ws)
}

// GethWs ethclient for websocket
func (nc *NodeErcClients) GethWs() *gethclient.Client {
	return gethclient.New(nc.Ws)
}
