package kecacc

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/loic-roux-404/crypto-bots/internal/model/token"
	"github.com/loic-roux-404/crypto-bots/internal/model/tx"
)

// ErrInvalidTx data
var (
	ErrInvalidTx = errors.New("Invalid transaction data")
	ErrEmptyRaw  = errors.New("Empty raw transaction")
	ErrNilErcTx  = errors.New("Empty erc transaction struct")
)

// NewTx prepare transaction requirements
func NewTx(
	to common.Address,
	nonce *big.Int,
	amount *big.Float,
	gasLimit *big.Int,
	gasPrice *big.Int,
	data []byte,
) (*tx.Adapter, error) {

	if data == nil || len(data) <= 0 {
		data = []byte{}
	}

	finalAmount := token.ToWei(amount)

	return &tx.Adapter{
		Hash:        "", // empty hash for non broadcasted tx
		To:          to.Hex(),
		Nonce:       nonce,
		TokenAmount: amount,
		Amount:      finalAmount,
		GasLimit:    gasLimit,
		GasPrice:    gasPrice,
		Data:        data,
	}, nil
}

// TxIsFrom address
func TxIsFrom(tx *types.Transaction, address common.Address) (bool, error) {
	if tx == nil {
		return false, ErrNilErcTx
	}

	err := IsErrAddress(address)
	if err != nil {
		return false, err
	}

	msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
	if err != nil {
		println(err.Error())
		return false, err
	}

	return msg.From() == address, nil
}

// KeccacTxFromRaw decode kecacc signer transaction hash string
func KeccacTxFromRaw(rawTx string) (*types.Transaction, error) {
	var err error
	if len(rawTx) <= 0 {
		return nil, ErrEmptyRaw
	}

	tx := new(types.Transaction)
	rawTxBytes, err := hex.DecodeString(rawTx)

	if err != nil {
		return nil, err
	}

	err = rlp.DecodeBytes(rawTxBytes, &tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// FromKeccacTx in go-etherum types.Transaction
func FromKeccacTx(keccacTx *types.Transaction) *tx.Adapter {
	return &tx.Adapter{
		Hash:        keccacTx.Hash().String(),
		To:          keccacTx.To().String(),
		Nonce:       new(big.Int).SetUint64(keccacTx.Nonce()),
		GasLimit:    nil,
		GasPrice:    keccacTx.GasPrice(),
		TokenAmount: token.FromWei(keccacTx.Value()),
		Amount:      keccacTx.Value(),
		Data:        keccacTx.Data(),
	}
}

// FromRaw decode kecacc signer transaction hash type
// rawTx need to be to rlp format
func FromRaw(rawTx string) (*tx.Adapter, error) {
	keccacTx, err := KeccacTxFromRaw(rawTx)

	if err != nil {
		return nil, err
	}

	return FromKeccacTx(keccacTx), nil
}

// ToRlp raw string
func ToRlp(keccacSignedTx *types.Transaction) (string, error) {
	rawTxBytes, err := keccacSignedTx.MarshalBinary()

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rawTxBytes), nil
}
