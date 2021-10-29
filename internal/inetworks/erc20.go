package inetworks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/model/account"
	"github.com/loic-roux-404/crypto-bots/internal/model/net"
	"github.com/loic-roux-404/crypto-bots/internal/model/token"
)

const (
	// ErcNetName identifier
	ErcNetName = "erc20"
	defaultNode = "ropsten"
)

// errors
var (
	errorImpossibleNonce = errors.New("Unable to determine nonce ")
)

// ErcHandler Handler config
type ErcHandler struct {
	name string
	client *ethclient.Client
	kecacc *account.Kecacc256
	config *net.ERCConfig
}

// NewEth create etherum handler
func NewEth() (Network, error) {
	cnf, err := net.NewERCConfig(ErcNetName, defaultNode); if err != nil {
		return nil, err
	}

	log.Printf("Connecting to %s...", cnf.Ipc)
	conn, err := ethclient.Dial(cnf.Ipc)

	if err != nil {
		return nil, err
	}

	acc, err := account.NewKecacc256(cnf.Pass, cnf.Keystore, cnf.FromAccount)

	if (err != nil) {
		return nil, err
	}

	return &ErcHandler{
		name: ErcNetName,
		kecacc: acc,
		client: conn,
		config: cnf,
	}, nil
}

// Send transaction to address
// Central function which need defer after the call
func (e *ErcHandler) Send(
	address string,
	amount *big.Float,
) (hash common.Hash, err error) {
	nonce, err := e.getNonce()
	// Create new transaction
	tx, err := e.createTx(address, nonce, amount, nil)

	if err != nil {
		return getEmptyHash(), err
	}

	return e.signAndBroadcast(tx)
}

// Update transaction
// Central function which need defer after the call
func (e *ErcHandler) Update(
	nonce *big.Int,
	address string,
	amount *big.Float,
) (hash common.Hash, err error) {
	// Create new transaction
	tx, err := e.createTx(address, nonce, amount, nil)

	if err != nil {
		return getEmptyHash(), err
	}

	return e.signAndBroadcast(tx)
}

// Cancel cancel transaction
func (e *ErcHandler) Cancel(nonce *big.Int) (common.Hash, error) {
	tx, err := e.createTx(e.kecacc.Account().Address.String(), nonce, big.NewFloat(0.0), nil)

	if err != nil {
		return getEmptyHash(), err
	}

	return e.signAndBroadcast(tx)
}

// TODO follow https://goethereumbook.org/address-check/

// Approve smart contract
// func (e *ErcHandler) Approve(address string) (hash common.Hash, err error) {
// 	// TODO address validator
// 	finalAddress := common.HexToAddress(address)
// 	return [20]byte{0}, nil
// }

// // Call smart contract method
// func (e *ErcHandler) Call(address string) (hash common.Hash, err error) {
// 	finalAddress := common.HexToAddress(address)
// 	return []byte{0}, nil
// }

func (e *ErcHandler) signAndBroadcast(tx *types.Transaction) (common.Hash, error) {
	signTx, err := e.signTx(tx)

	if err != nil {
		return getEmptyHash(), err
	}

	return e.broadcastTx(signTx)
}

// Broadcast transaction to network
func (e *ErcHandler) broadcastTx(signTx *types.Transaction) (common.Hash, error) {
	// Send the transaction
	log.Printf("Sending transaction on chain : %d", e.config.ChainID)
	// TODO use goroutine
	err := e.client.SendTransaction(context.Background(), signTx)

	if err != nil {
		return getEmptyHash(), err
	}

	// Obtain transaction hash as a string
	return signTx.Hash(), nil
}

// Sign this transaction with current account
func (e *ErcHandler) signTx(tx *types.Transaction) (*types.Transaction, error) {
	// Sign the transaction with private key
	signTx, err := e.kecacc.Ks().SignTx(
		e.kecacc.Account(),
		tx,
		big.NewInt(e.config.ChainID),
	)

	if err != nil {
		return nil, err
	}

	return signTx, nil
}

func (e *ErcHandler) setFees(address common.Address) error {
	if e.config.ManualFee {
		return nil
	}

	limit, err := e.client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &address,
		Data: []byte{0},
	})

	if err != nil {
		return err
	}

	e.config.GasLimit = limit

	price, err := e.client.SuggestGasPrice(context.Background())

	if err != nil {
		return err
	}

	e.config.GasPrice = price.Int64()

	return nil
}

func (e *ErcHandler) getNonce() (*big.Int, error) {
	// NonceAt returns the account nonce of the given account.
	// nonce, err := client.NonceAt(ctx, address, nil)
	a := e.kecacc.Account().Address
	// This is the nonce that should be used for the next transaction.
	nonce, err := e.client.PendingNonceAt(context.Background(), a)
	finalNonce := new(big.Int).SetUint64(nonce)

	if err != nil {
		return nil, fmt.Errorf("%s: %s", errorImpossibleNonce, err)
	}

	return finalNonce, nil
}

func (e *ErcHandler) createTx(
	address string,
	nonce *big.Int,
	amount *big.Float,
	data []byte,
) (*types.Transaction, error) {
		// prepare transaction requirements
		panicIfInvalidAddress(address)
		finalAddress := common.HexToAddress(address)
		e.setFees(finalAddress)

		if (data == nil || len(data) <= 0) {
			data = []byte{}
		}

		finalAmount := token.EtherToWei(amount)

		logTx(helpers.Map{
			"nonce": nonce,
			"from": e.kecacc.Account().Address,
			"to": finalAddress,
			"data": data,
			"gasLimit": e.config.GasLimit,
			"gasPrice": e.config.GasPrice,
			"Wei": finalAmount,
			"Eth": amount,
		})

		// Create new transaction
		tx := types.NewTransaction(
			nonce.Uint64(),
			finalAddress,
			finalAmount,
			e.config.GasLimit,
			big.NewInt(e.config.GasPrice),
			data,
		)

		return tx, nil
}

func panicIfInvalidAddress(address string) {
	// prepare transaction requirements
	acc := accounts.Account{
		Address: common.HexToAddress(address),
	}
	isValidAd := account.ValidateAddress(acc)

	if !isValidAd {
		log.Panic(account.ErrInvalid.Error())
	}
}

func logTx(m helpers.Map) {
	jsonString, _ := json.Marshal(m)

	log.Printf("info: Tx %s", jsonString)
}

func getEmptyHash() common.Hash {
	return (common.Hash)([common.HashLength]byte{0})
}
