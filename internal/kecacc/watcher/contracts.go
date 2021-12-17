package watcher

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"

	"github.com/loic-roux-404/crypto-bots/internal/model/sub"
)

// Sc type
type Sc struct {
	wsClient   *ethclient.Client
	gethClient *gethclient.Client
	Sub        ethereum.Subscription
	logs       chan types.Log
}

// NewSc for a specific query (need websocket connection)
// Watch smart contract events
func NewSc(ws *ethclient.Client, q ethereum.FilterQuery) (*Sc, error) {
	w := &Sc{wsClient: ws}
	sub, logs, err := w.sub(q)

	if err != nil {
		return nil, err
	}

	w.logs = logs
	w.Sub = sub

	return w, nil
}

// RunEventLoop launch
func (w *Sc) RunEventLoop(callback sub.ScSubCallback) {
	for {
		select {
		case err := <-w.Sub.Err():
			log.Fatal(err)
		case vLog := <-w.logs:
			fmt.Println(vLog)
			callback(vLog)
		}
	}
}

// Sub builder
func (w *Sc) sub(q ethereum.FilterQuery) (ethereum.Subscription, chan types.Log, error) {
	logs := make(chan types.Log)
	// TODO if smart contract subscribe filter SubscribeFilterLogs
	sub, err := w.wsClient.SubscribeFilterLogs(context.Background(), q, logs)
	if err != nil {
		return nil, nil, err
	}

	return sub, logs, nil
}
