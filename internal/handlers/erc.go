package handlers

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
)

type ercH struct {
	client *ethclient.Client
}

func (e *ercH) NewEth() (*ercH, error) {
	conn, err := ethclient.Dial(os.Getenv("ERC_IPC"),)

	if err != nil {
	  return nil, fmt.Errorf("Failed to connect to the Ethereum client: %v", err)
	}

	return &ercH{client: conn}, nil
}

func (e *ercH) Send(address string) {
	// Create new transaction
	
	tx := types.NewTransaction(
		e.client.getNonce(),
		toAddress,
		amount,
		gasLimit,
		gasPrice,
		data,
	)

	// Sign the transaction with private key
	signTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey) // See other signers in transaction_signing.go file in go-ethereum project

	if err != nil {
		panic(err)
	}

	// Send the transaction
	err = e.client.SendTransaction(ctx, signTx)

	if err != nil {
		panic(err)
	}
	
	// Obtain transaction hash as a string
	strHash := signTx.Hash().String()
}

func (e *ercH) getNonce(address string) {
	// NonceAt returns the account nonce of the given account.
	// nonce, err := client.NonceAt(ctx, address, nil)

	// This is the nonce that should be used for the next transaction.
	nonce, err := e.client.PendingNonceAt(ctx, address)

	if err != nil {
		return nil, fmt.Errorf("Unable to determine nonce")
	}

	return nonce, nil
}


func (e *ercH) parseAddress(address string) {
	a := common.HexToAddress(address)
	return a
}