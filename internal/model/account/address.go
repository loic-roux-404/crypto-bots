package account

import (
	"context"
	"errors"
	"regexp"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ErrInvalid invalid
var ErrInvalid = errors.New("Invalid address")

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
		nil,// nil is latest block
	)

	if err != nil {
	  return false, err
	}

	return len(bytecode) > 0, nil
}
