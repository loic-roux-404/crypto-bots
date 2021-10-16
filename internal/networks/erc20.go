package networks

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Config of etherul handler
type config struct {
	gasLimit *big.Int
	gasPrice *big.Int 
}

// NewConf of erc handler
func newConf(gasLimit *big.Int, gasPrice *big.Int) (*config)  {
	cnf := &config{big.NewInt(91.00), big.NewInt(9.00)}
	// TODO get consistent gas price
	if (gasLimit != nil) {
		cnf.gasLimit = gasPrice
	}

	if (gasPrice != nil) {
		cnf.gasPrice = gasPrice
	}

	return cnf
}

// ErcHandler Handler config
type ErcHandler struct {
	client *ethclient.Client
	currentPrivateKey *ecdsa.PrivateKey
	config *config
}

// NewEth create etherum handler
func NewEth() (Network, error) {
	conn, err := ethclient.Dial(os.Getenv("ERC_IPC"),)

	if err != nil {
	  return nil, fmt.Errorf("Failed to connect to the Ethereum client: %v", err)
	}

	return &ErcHandler{client: conn, config: newConf(nil, nil)}, nil
}

// Send transaction to address
func (e *ErcHandler) Send(address string, amount *big.Int) (hash common.Hash, err error) {
	// prepare transaction requirements
	finalAddress, err := e.parseAddress("") // TODO

	if err != nil {
		panic(err)
	}

	nonce, err := e.getNonce(finalAddress)

	if err != nil {
		panic(err)
	}

	data := []byte{}
	// Create new transaction
	tx := types.NewTransaction(
		nonce.Uint64(),
		finalAddress, 
		amount,
		e.config.gasLimit.Uint64(),
		e.config.gasPrice,
		data,
	)
	
	// Sign the transaction with private key
	signTx, err := types.SignTx(tx, types.HomesteadSigner{}, e.initPrivateKeyFromMem("")) // See other signers in transaction_signing.go file in go-ethereum project

	if err != nil {
		panic(err)
	}

	// Send the transaction
	err = e.client.SendTransaction(context.Background(), signTx)

	if err != nil {
		panic(err)
	}
	
	// Obtain transaction hash as a string
	return signTx.Hash(), nil
}

// Call smart contract method
func (e *ErcHandler) Call(address string) (hash *common.Hash, err error) {
	return nil, nil
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

func (e *ErcHandler) initPrivateKeyFromMem(key string) *ecdsa.PrivateKey {
	// TODO memstorage call
	return nil
}

func (e *ErcHandler) parseAddress(address string) (common.Address, error) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	if (!re.MatchString(address)) {
		return common.HexToAddress(address), fmt.Errorf("Address is invalid: %v", address) 
	}

	return common.HexToAddress(address), nil
}
