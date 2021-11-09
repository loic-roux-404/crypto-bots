package erc20

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/model/kecacc"
	"github.com/loic-roux-404/crypto-bots/internal/model/net"
	"github.com/loic-roux-404/crypto-bots/internal/model/store"
	"github.com/loic-roux-404/crypto-bots/internal/model/token"
	"github.com/loic-roux-404/crypto-bots/internal/watcher"
)

const (
	// ErcNetName identifier
	ErcNetName  = "erc20"
	defaultNode = "ropsten"
)

// errors
var (
	errorImpossibleNonce = errors.New("Unable to determine nonce ")
)

// ErcHandler Handler config
type ErcHandler struct {
	name      string
	clients   *NodeClients
	kecacc    *kecacc.KeccacWallet
	config    *net.ERCConfig
	contracts map[string]*store.Contract
}

// NewEth create etherum handler
func NewEth() net.Network {
	cnf, err := net.NewERCConfig(ErcNetName, defaultNode)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Info: Connecting to %s...", cnf.Ipc)
	clients, err := NewClients(cnf.Ipc, cnf.Ws)

	if err != nil {
		log.Fatal(err)
	}

	acc, err := kecacc.NewErcWallet(cnf.Pass, cnf.Keystore, cnf.FromAccount)

	if err != nil {
		log.Fatal(err)
	}

	return &ErcHandler{
		name:      ErcNetName,
		kecacc:    acc,
		clients:   clients,
		config:    cnf,
		contracts: make(map[string]*store.Contract),
	}
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
func (e *ErcHandler) Cancel(nonce *big.Int) common.Hash {
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
func (e *ErcHandler) Deploy(input string, storeDeployFn store.DeployFn) interface{} {
	defer helpers.RecoverAndLog()
	auth, err := e.getAuth()

	if err != nil {
		panic(err)
	}

	address, tx, instance, err := storeDeployFn(auth, e.clients.EthRPC(), input)
	addressStr := address.String()
	e.contracts[addressStr] = store.NewContract(address, tx, instance)

	if err != nil {
		panic(err)
	}

	log.Printf("Contract deployed: %s", e.contracts[addressStr].JSON())

	return e.contracts[addressStr]
}

// Load a smart contract
func (e *ErcHandler) Load(address string, loadFn store.LoadFn) interface{} {
	defer helpers.RecoverAndLog()
	finalAddress := common.HexToAddress(address)
	isScAddress, err := kecacc.ValidateSc(e.clients.EthRPC(), common.HexToAddress(address))

	if value, ok := e.contracts[address]; ok {
		return value
	}

	if !isScAddress || err != nil {
		panic(fmt.Errorf(store.ErrLoadSc.Error(), address, err))
	}

	instance, err := loadFn(finalAddress, e.clients.EthRPC())
	e.contracts[address] = store.NewContract(finalAddress, nil, instance)

	if err != nil {
		panic(err)
	}

	return e.contracts[address]
}

// Subscribe an address
// TODO move all in watcher module
func (e *ErcHandler) Subscribe(address string) watcher.WatcherSc {
	log.Printf("Info: Booting subscriber on: %s", address)

	finalAddress := common.HexToAddress(address)

	isSc, err := kecacc.ValidateSc(e.clients.EthRPC(), finalAddress)
	if err != nil {
		panic(err)
	}

	if !isSc {
		panic(kecacc.ErrScInvalid)
	}

	w, err := e.subscribeSc(finalAddress)

	if err != nil {
		panic(err)
	}

	return w

}

// SubscribeCurrent account
func (e *ErcHandler) SubscribeCurrent() watcher.WatcherAcc {

	w, err := e.subscribeAcc(e.kecacc.Account().Address)

	if err != nil {
		panic(err)
	}

	return w
}

func (e *ErcHandler) subscribeSc(address common.Address) (w *watcher.Sc, err error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
	}

	w, err = watcher.NewSc(e.clients.EthWs(), query)

	if err != nil {
		return nil, err
	}

	return w, nil
}

func (e *ErcHandler) subscribeAcc(address common.Address) (w *watcher.Acc, err error) {
	w, err = watcher.NewAcc(e.clients.GethWs(), e.kecacc)

	if err != nil {
		return nil, err
	}

	return w, nil
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
	err := e.clients.EthRPC().SendTransaction(context.Background(), signTx)

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

	limit, err := e.clients.EthRPC().EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &address,
		Data: []byte{0},
	})

	if err != nil {
		return err
	}

	e.config.GasLimit = limit

	price, err := e.clients.EthRPC().SuggestGasPrice(context.Background())

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
	nonce, err := e.clients.EthRPC().PendingNonceAt(context.Background(), a)
	finalNonce := new(big.Int).SetUint64(nonce)

	if err != nil {
		return nil, fmt.Errorf("Error: %s %s", errorImpossibleNonce, err)
	}

	return finalNonce, nil
}

// prepare transaction requirements
// TODO refacto to transaction module (named keccac tx)
func (e *ErcHandler) createTx(
	address string,
	nonce *big.Int,
	amount *big.Float,
	data []byte,
) (*types.Transaction, error) {
	err := kecacc.IsErrStrAddress(address)

	if err != nil {
		return nil, err
	}

	finalAddress := common.HexToAddress(address)
	err = e.setFees(finalAddress)

	if err != nil {
		return nil, err
	}

	tx, err := kecacc.NewTx(
		finalAddress,
		nonce,
		amount,
		new(big.Int).SetUint64(e.config.GasLimit),
		big.NewInt(e.config.GasPrice),
		data,
	)

	if err != nil {
		return nil, err
	}

	// Create new transaction
	ercTx := types.NewTransaction(
		tx.Nonce.Uint64(),
		tx.To,
		tx.Amount,
		tx.GasLimit.Uint64(),
		tx.GasPrice,
		tx.Data,
	)

	tx.Hash = ercTx.Hash()
	tx.Log()

	return ercTx, nil
}

func (e *ErcHandler) getAuth() (*bind.TransactOpts, error) {
	acc := e.kecacc.Account()
	auth, err := bind.NewKeyStoreTransactor(e.kecacc.Ks(), acc)

	if err != nil {
		return nil, err
	}

	auth.Nonce, err = e.getNonce()
	if err != nil {
		return nil, err
	}

	err = e.setFees(acc.Address)
	if err != nil {
		return nil, err
	}

	auth.Value = token.ToWei(big.NewFloat(0.00))
	auth.GasLimit = e.config.GasLimit
	auth.GasPrice = big.NewInt(e.config.GasPrice)

	return auth, nil
}
