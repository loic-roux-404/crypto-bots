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
	"github.com/loic-roux-404/crypto-bots/internal/kecacc"
	"github.com/loic-roux-404/crypto-bots/internal/kecacc/fees"
	"github.com/loic-roux-404/crypto-bots/internal/kecacc/store"
	"github.com/loic-roux-404/crypto-bots/internal/kecacc/watcher"
	"github.com/loic-roux-404/crypto-bots/internal/model/net"
	"github.com/loic-roux-404/crypto-bots/internal/model/sub"
	"github.com/loic-roux-404/crypto-bots/internal/nets/erc20/clients"
)

const (
	// ErcNetName identifier
	ErcNetName = "erc20"
	DefaultNet = "ropsten"
)

// errors
var (
	ErrImpossibleNonce = errors.New("Unable to determine nonce ")
)

// ErcHandler Handler config
type ErcHandler struct {
	Name      string
	clients   *clients.NodeErcClients
	kecacc    *kecacc.KeccacWallet
	config    *net.Config
	contracts map[string]*store.Contract
}

// NewEth create etherum handler
func NewEth(cnf *net.Config) net.Network {

	log.Printf("Info: Connecting to %s...", cnf.Ipc)
	clients, err := clients.NewClients(cnf.Ipc, cnf.Ws)

	if err != nil {
		log.Fatal(err)
	}

	acc, err := kecacc.NewWallet(cnf.Pass, cnf.Keystore, cnf.FromAddress, cnf.Wallets)

	if err != nil {
		log.Fatal(err)
	}

	return &ErcHandler{
		Name:      ErcNetName,
		kecacc:    acc,
		clients:   clients,
		config:    cnf,
		contracts: make(map[string]*store.Contract),
	}
}

// Send transaction to address
// Central function which need defer after the call
func (e *ErcHandler) Send(address string, amount *big.Float) string {
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

	return sentTx.String()
}

// Update transaction
// Central function which need defer after the call
func (e *ErcHandler) Update(
	nonce *big.Int,
	address string,
	amount *big.Float,
) string {
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

	return sentTx.String()
}

// Cancel cancel transaction
func (e *ErcHandler) Cancel(nonce *big.Int) string {
	defer helpers.RecoverAndLog()
	tx, err := e.createTx(e.kecacc.Account().Address.String(), nonce, big.NewFloat(0.0), nil)

	if err != nil {
		panic(err)
	}

	sentTx, err := e.signAndBroadcast(tx)
	if err != nil {
		panic(err)
	}

	return sentTx.String()
}

// CurrentBalance of account logged in
func (e *ErcHandler) CurrentBalance() *big.Float {
	defer helpers.RecoverAndLog()

	b, err := e.balanceAt(e.kecacc.Account().Address.String())

	if err != nil {
		panic(err)
	}

	return b
}

// BalanceAt Logged in Account balance
func (e *ErcHandler) BalanceAt(address string) *big.Float {
	defer helpers.RecoverAndLog()
	b, err := e.balanceAt(address)

	if err != nil {
		panic(err)
	}

	return b
}

// balanceAt generic eth client call to fetch balance
func (e *ErcHandler) balanceAt(address string) (*big.Float, error) {
	if kecacc.ValidateAddress(address) {
		return nil, kecacc.ErrAddressInvalid
	}

	b, err := e.clients.EthRPC().BalanceAt(
		context.Background(),
		common.HexToAddress(address),
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("error: impossible to retrieve current account balance : %s", err)
	}

	return fees.WeiToDecimal(b), err
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
func (e *ErcHandler) Subscribe(address string) sub.Sc {
	log.Printf("Info: Booting subscriber on: %s", address)

	finalAddress := common.HexToAddress(address)

	_, err := kecacc.ValidateSc(e.clients.EthRPC(), finalAddress)
	if err != nil {
		panic(err)
	}

	w, err := watcher.NewSc(e.clients.EthWs(), ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(address)},
	})

	if err != nil {
		panic(err)
	}

	return w
}

// SubscribeCurrent account
func (e *ErcHandler) SubscribeCurrent() sub.Acc {
	w, err := watcher.NewAcc(e.clients, e.kecacc)

	if err != nil {
		panic(err)
	}

	return w
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
		return nil, fmt.Errorf("Error: %s %s", ErrImpossibleNonce, err)
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

	log.Println(kecacc.ToRlp(tx.Adapted.(*types.Transaction)))

	tx.Log()

	return tx.Adapted.(*types.Transaction), nil
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

	auth.Value = fees.ToWei(big.NewFloat(0.00))
	auth.GasLimit = e.config.GasLimit
	auth.GasPrice = big.NewInt(e.config.GasPrice)

	return auth, nil
}
