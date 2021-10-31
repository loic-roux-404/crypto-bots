package transaction

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/model/token"
)

// Tx transaction infos
// TODO more generic address (non keccac related)
type Tx struct {
	Hash common.Hash
	To common.Address
	Nonce *big.Int
	GasLimit *big.Int
	GasPrice *big.Int
	TokenAmount *big.Float
	Amount *big.Int
	Data []byte
}

// NewTx prepare transaction requirements
func NewTx(
	to common.Address,
	nonce *big.Int,
	amount *big.Float,
	gasLimit *big.Int,
	gasPrice *big.Int,
	data []byte,
) (*Tx, error) {

		if (data == nil || len(data) <= 0) {
			data = []byte{}
		}

		finalAmount := token.ToWei(amount)

		return &Tx{
			Hash: common.Hash{}, // empty hash for non broadcasted tx
			To: to,
			Nonce: nonce,
			TokenAmount: amount,
			Amount: finalAmount,
			GasLimit: gasLimit,
			GasPrice: gasPrice,
			Data: data,
		}, nil
}

// TxIsFrom address
// TODO validate address / tx hash
func TxIsFrom(hash common.Hash, address common.Address) (bool, error) {

	publicKeyBytes, err := hex.DecodeString(hash.String())

	if err != nil {
		return false, err
	}

	sigPublicKey, err := crypto.Ecrecover(address[:], []byte{})

	if err != nil {
		return false, err
	}

	return bytes.Equal(sigPublicKey, publicKeyBytes), nil
}

// KeccacTx decode kecacc signer transaction hash string
func KeccacTx(rawTx common.Hash) (tx *types.Transaction, err error) {
    rawTxBytes, err := hex.DecodeString(rawTx.String())
	rlp.DecodeBytes(rawTxBytes, &tx)

	return tx, err
}

// NewTxFromKeccacHash decode kecacc signer transaction hash type
func NewTxFromKeccacHash(hash common.Hash) (*Tx, error) {
	keccacTx, err := KeccacTx(hash)

	if err != nil {
		return nil, err
	}

	return &Tx{
		Hash: hash,
		To: *keccacTx.To(),
		Nonce: new(big.Int).SetUint64(keccacTx.Nonce()),
		GasLimit: nil,
		GasPrice: keccacTx.GasPrice(),
		TokenAmount: token.FromWei(keccacTx.Value()),
		Amount: keccacTx.Value(),
		Data: keccacTx.Data(),
	}, nil
}

// Log transaction in json
func (tx *Tx) Log() {
	m := helpers.Map{
		"nonce": tx.Nonce.Uint64(),
		"to": tx.To,
		"data": tx.Data,
		"gasLimit": tx.GasLimit,
		"gasPrice": tx.GasPrice,
		"Wei": tx.Amount,
		"Eth": tx.TokenAmount,
	}
	jsonString, _ := json.Marshal(m)

	log.Printf("info: Tx %s", jsonString)
}
