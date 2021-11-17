package watcher

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/loic-roux-404/crypto-bots/internal/kecacc"
	"github.com/loic-roux-404/crypto-bots/internal/model/sub"
	"github.com/loic-roux-404/crypto-bots/internal/nets/erc20/clients"
)

// Acc type
type Acc struct {
	clients    *clients.NodeErcClients
	Sub        ethereum.Subscription
	txs        chan common.Hash
	acc        *kecacc.KeccacWallet
}

// NewAcc for a specific query (need websocket connection)
func NewAcc(clients *clients.NodeErcClients, acc *kecacc.KeccacWallet) (w *Acc, err error) {
	w = &Acc{clients: clients, acc: acc}
	sub, txs, err := w.sub()

	if err != nil {
		return nil, err
	}

	w.Sub = sub
	w.txs = txs

	return w, nil
}

// RunEventLoop launch
func (w *Acc) RunEventLoop(callback sub.AccSubCallback) {
	for {
		select {
		case err := <-w.Sub.Err():
			log.Fatal(err)
		case vTx := <-w.txs:
			tx, _, _ := w.clients.EthRPC().TransactionByHash(context.Background(), vTx)

			ok, err := w.acc.IsTxFromCurrent(tx); if err != nil {
				log.Panic(err)
			}

			if !ok {
				continue
			}

			finalTx := kecacc.FromKeccacTx(tx)

			if err != nil {
				log.Panic(err)
			}

			log.Printf("Info: callback on : %s", vTx)
			callback(finalTx)
		}
	}
}

// Sub build subscriber
func (w *Acc) sub() (ethereum.Subscription, chan common.Hash, error) {

	txs := make(chan common.Hash)
	sub, err := w.clients.GethWs().SubscribePendingTransactions(context.Background(), txs)

	if err != nil {
		return nil, nil, err
	}

	return sub, txs, nil
}
