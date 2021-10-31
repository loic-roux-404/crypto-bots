package inetworks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/model/account"
	"github.com/loic-roux-404/crypto-bots/internal/model/net"
	"github.com/loic-roux-404/crypto-bots/internal/model/token"
	"github.com/loic-roux-404/crypto-bots/internal/model/store"
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
	contracts map[string]*store.Contract
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
		contracts: make(map[string]*store.Contract),
	}, nil
}

// Send transaction to address
// Central function which need defer after the call
func (e *ErcHandler) Send(
	address string,
	amount *big.Float,
) (hash common.Hash) {
	defer helpers.RecoverAndLog()
	nonce, err := e.getNonce()

	if err != nil {
		panic(err)
	}

	// Create new transaction
	tx, err := e.createTx(address, nonce, amount, nil)

	if err != nil {
		panic(err)
	}

	sentTx, err := e.signAndBroadcast(tx)

	if err != nil {
		panic(err)
	}

	return sentTx
}

// Update transaction
// Central function which need defer after the call
func (e *ErcHandler) Update(
	nonce *big.Int,
	address string,
	amount *big.Float,
) (hash common.Hash) {
	defer helpers.RecoverAndLog()
	// Create new transaction
	tx, err := e.createTx(address, nonce, amount, nil)

	if err != nil {
		panic(err)
	}

	sentTx, err := e.signAndBroadcast(tx)

	if err != nil {
		panic(err)
	}

	return sentTx
}

// Cancel cancel transaction
func (e *ErcHandler) Cancel(nonce *big.Int) (common.Hash) {
	defer helpers.RecoverAndLog()
	tx, err := e.createTx(e.kecacc.Account().Address.String(), nonce, big.NewFloat(0.0), nil)

	if err != nil {
		panic(err)
	}

	sentTx, err := e.signAndBroadcast(tx)

	if err != nil {
		panic(err)
	}

	return sentTx
}

// Deploy a smart contract api function
func (e *ErcHandler) Deploy(
	input string,
	storeDeployFn store.DeployFn,
) *store.Contract {
	defer helpers.RecoverAndLog()
	auth, err := e.getAuth()

	if err != nil {
		panic(err)
	}

	address, tx, instance, err := storeDeployFn(auth, e.client, input)
	addressStr := address.String()
	e.contracts[addressStr] = store.NewContract(address, tx, instance)

	if err != nil {
		panic(err)
	}

	log.Printf("Contract deployed: %s", e.contracts[addressStr].JSON())

	return e.contracts[addressStr]
}

// Load a smart contract
func (e *ErcHandler) Load(address string, loadFn store.LoadFn) *store.Contract {
	defer helpers.RecoverAndLog()
	finalAddress := common.HexToAddress(address)
	isScAddress, err := account.ValidateSc(e.client, common.HexToAddress(address))

	if value, ok := e.contracts[address]; ok {
		return value
	}

	if !isScAddress || err != nil {
		panic(fmt.Errorf(store.ErrLoadSc.Error(), address, err))
	}

	instance, err := loadFn(finalAddress, e.client)
	e.contracts[address] = store.NewContract(finalAddress, nil, instance)

	if err != nil {
        panic(err)
    }

    return e.contracts[address]
}

func (e *ErcHandler) signAndBroadcast(tx *types.Transaction) (common.Hash, error) {
	signTx, err := e.signTx(tx)

	if err != nil {
		return common.Hash{}, err
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
		return common.Hash{}, err
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
		return nil, fmt.Errorf("Error: %s %s", errorImpossibleNonce, err)
	}

	return finalNonce, nil
}

// prepare transaction requirements
// TODO : move it in model/transaction
func (e *ErcHandler) createTx(
	address string,
	nonce *big.Int,
	amount *big.Float,
	data []byte,
) (*types.Transaction, error) {
		err := account.IsErrAddress(address)

		if err != nil {
			return nil, err
		}

		finalAddress := common.HexToAddress(address)
		err = e.setFees(finalAddress)

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
			e.config.GasLimit,
			big.NewInt(e.config.GasPrice),
			data,
		)

		return tx, nil
}

func (e *ErcHandler) getAuth() (*bind.TransactOpts, error) {
	acc := e.kecacc.Account()
	auth, err := bind.NewKeyStoreTransactor(e.kecacc.Ks(), acc)

	if err != nil {
		return nil, err
	}

	auth.Nonce, err = e.getNonce(); if err != nil {
		return nil, err
	}

	err = e.setFees(acc.Address); if err != nil {
		return nil, err
	}

	auth.Value = token.EtherToWei(big.NewFloat(0.00))
	auth.GasLimit = e.config.GasLimit
	auth.GasPrice = big.NewInt(e.config.GasPrice)

	return auth, nil
}

// TODO move in model/transaction
func logTx(m helpers.Map) {
	jsonString, _ := json.Marshal(m)

	log.Printf("info: Tx %s", jsonString)
}

func getEmptyHash() common.Hash {
	return (common.Hash)([common.HashLength]byte{0})
}
