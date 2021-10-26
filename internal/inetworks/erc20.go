package inetworks

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

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

	acc, err := account.NewKecacc256(cnf.Memonic, cnf.Keystore)

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
	pair token.Pair,
	amount *big.Int,
) (hash common.Hash, err error) {
	// Create new transaction
	tx, err := e.createTx(address, amount, nil)

	// Sign the transaction with private key
	signTx, err := e.kecacc.Store().SignTx(
		e.kecacc.Account(),
		tx,
		big.NewInt(3),
	)

	if err != nil {
		log.Panic(err)
	}

	// Send the transaction
	err = e.client.SendTransaction(context.Background(), signTx)

	if err != nil {
		log.Panic(err)
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
	estimatedGas, err := e.client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &address,
		Data: []byte{0},
	})

	if err != nil {
		return err
	}

	gasLimit := int64(float64(estimatedGas) * 1.30)
	fmt.Println(gasLimit) // 27305

	e.config.GasLimit = gasLimit

	return nil
}

func (e *ErcHandler) getNonce(address common.Address) (*big.Int, error) {
	// NonceAt returns the account nonce of the given account.
	// nonce, err := client.NonceAt(ctx, address, nil)

	// This is the nonce that should be used for the next transaction.
	nonce, err := e.client.PendingNonceAt(context.Background(), address)
	finalNonce := new(big.Int).SetUint64(nonce)

	if err != nil {
		return nil, fmt.Errorf("Unable to determine nonce : %s", err)
	}

	return finalNonce, nil
}

func (e *ErcHandler) createTx(
	address string,
	amount *big.Int,
	data []byte,
) (*types.Transaction, error) {
		// prepare transaction requirements
		finalAddress := common.HexToAddress(address)
		e.estimateGas(finalAddress)

		nonce, err := e.getNonce(finalAddress)

		if err != nil {
			return nil, err
		}

		if (data == nil || len(data) > 0) {
			data = []byte{}
		}

		// Create new transaction
		tx := types.NewTransaction(
			nonce.Uint64(),
			finalAddress,
			amount,
			big.NewInt(e.config.GasLimit).Uint64(),
			big.NewInt(e.config.GasPrice),
			data,
		)

		return tx, nil
}

func valAndGetAddress(address string) common.Address {
	// prepare transaction requirements
	finalAddress := common.HexToAddress(address)

	return finalAddress
}
