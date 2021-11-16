package tx

import (
	"encoding/json"
	"log"
	"math/big"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
)

// Adapter transaction infos
// TODO more generic address (non keccac related)
type Adapter struct {
	Hash        string
	To          string
	Nonce       *big.Int
	GasLimit    *big.Int
	GasPrice    *big.Int
	TokenAmount *big.Float
	Amount      *big.Int
	Data        []byte
}

// Log transaction in json
func (tx *Adapter) Log() {
	m := helpers.Map{
		"hash":     tx.Hash,
		"nonce":    tx.Nonce.Uint64(),
		"to":       tx.To,
		"data":     tx.Data,
		"gasLimit": tx.GasLimit,
		"gasPrice": tx.GasPrice,
		"Wei":      tx.Amount,
		"Eth":      tx.TokenAmount,
	}
	jsonString, _ := json.Marshal(m)

	log.Printf("info: Tx %s", jsonString)
}
