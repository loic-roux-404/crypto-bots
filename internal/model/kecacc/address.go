package kecacc

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ErrAccInvalid invalid
var (
	ErrAccInvalid = errors.New("Invalid public key address")
	ErrScInvalid  = errors.New("Invalid smart contract address")
)

// ValidateAddress destination
func ValidateAddress(acc accounts.Account) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	return re.MatchString(acc.Address.String())
}

// ValidateSc check if it's a smart contract
func ValidateSc(client *ethclient.Client, address common.Address) (bool, error) {
	bytecode, err := client.CodeAt(
		context.Background(),
		address,
		nil, // nil is latest block
	)

	if err != nil {
		return false, fmt.Errorf("%s : %s", ErrScInvalid, err)
	}

	return len(bytecode) > 0, nil
}

// ValidateTx from tx struct
// TODO nil tx error
func ValidateTx(tx *types.Transaction) bool {
	if tx == nil {
		return false
	}

	_, err := hexutil.Decode(tx.Hash().String())
	println("=====",err)
	return err == nil
}

// IsErrAddress validate address but return an error if invalid
func IsErrAddress(address common.Address) error {
	acc := accounts.Account{
		Address: address,
	}
	isValidAd := ValidateAddress(acc)

	if isValidAd {
		return nil
	}

	return fmt.Errorf("%s : %s", ErrAccInvalid, address)
}

// IsErrStrAddress is this string address ok
func IsErrStrAddress(address string) error {
	finalAddress := common.HexToAddress(address)

	return IsErrAddress(finalAddress)
}
