package inetworks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
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
		return nil, fmt.Errorf("Failed to connect to the Ethereum client: %v", err)
	}

	acc, err := account.NewKecacc256(cnf.Pass, cnf.Keystore)

	if (err != nil) {
		return nil, fmt.Errorf("Failed to init wallet: %v", err)
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
	// Create new transaction
	tx, err := e.createTx(address, amount, nil)

	// Sign the transaction with private key
	signTx, err := e.kecacc.Store().SignTx(
		e.kecacc.Account(),
		tx,
		big.NewInt(e.config.ChainID),
	)

	if err != nil {
		return (common.Hash)([common.HashLength]byte{0}), err
	}

	// Send the transaction
	log.Printf("Sending transaction on chain : %d", e.config.ChainID)
	err = e.client.SendTransaction(context.Background(), signTx)

	if err != nil {
		return (common.Hash)([common.HashLength]byte{0}), err
	}

	// Obtain transaction hash as a string
	return signTx.Hash(), nil
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

func (e *ErcHandler) estimateGas(address common.Address) error {
	if e.config.ManualFee {
		return nil
	}

	estimatedGas, err := e.client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &address,
		Data: []byte{0},
	})

	if err != nil {
		return err
	}

	e.config.GasLimit = int64(estimatedGas)

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
		return nil, fmt.Errorf("Unable to determine nonce : %s", err)
	}

	return finalNonce, nil
}

func (e *ErcHandler) createTx(
	address string,
	amount *big.Float,
	data []byte,
) (*types.Transaction, error) {
		// prepare transaction requirements
		finalAddress := common.HexToAddress(address)
		e.estimateGas(finalAddress)

		nonce, err := e.getNonce()

		if err != nil {
			return nil, err
		}

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
			uint64(e.config.GasLimit),
			big.NewInt(e.config.GasPrice),
			data,
		)

		return tx, nil
}

// TODO add checking regex for smart contract or address
func valAndGetAddress(address string) common.Address {
	// prepare transaction requirements
	finalAddress := common.HexToAddress(address)

	return finalAddress
}

func logTx(m helpers.Map) {
	jsonString, _ := json.Marshal(m)

	log.Printf("info: Tx %s", jsonString)
}
