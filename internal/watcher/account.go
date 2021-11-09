package watcher

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"

	"github.com/loic-roux-404/crypto-bots/internal/model/kecacc"
)

// Acc type
type Acc struct {
	gethClient *gethclient.Client
	Sub        ethereum.Subscription
	txs        chan common.Hash
	acc        *kecacc.KeccacWallet
}

// NewAcc for a specific query (need websocket connection)
func NewAcc(ws *gethclient.Client, acc *kecacc.KeccacWallet) (w *Acc, err error) {
	w = &Acc{gethClient: ws, acc: acc}
	sub, txs, err := w.sub()

	if err != nil {
		return nil, err
	}

	w.Sub = sub
	w.txs = txs

	return w, nil
}

// RunEventLoop launch
func (w *Acc) RunEventLoop(callback AccSubCallback) {
	for {
		select {
		case err := <-w.Sub.Err():
			log.Fatal(err)
		case vTx := <-w.txs:
			ok, err := w.acc.IsTxFromCurrent(vTx)
			if err != nil {
				log.Panic(err)
			}

			if !ok {
				return
			}

			tx, err := kecacc.NewTxFromKeccacHash(vTx)
			if err != nil {
				log.Panic(err)
			}
			log.Printf("Info: callback on : %s", tx.Hash)
			callback(tx)
		}
	}
}

// Sub build subscriber
func (w *Acc) sub() (ethereum.Subscription, chan common.Hash, error) {

	txs := make(chan common.Hash)
	sub, err := w.gethClient.SubscribePendingTransactions(context.Background(), txs)

	if err != nil {
		return nil, nil, err
	}

	return sub, txs, nil
}
