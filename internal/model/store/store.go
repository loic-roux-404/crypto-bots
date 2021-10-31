package store

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/loic-roux-404/crypto-bots/internal/helpers"
)

// Contract type
type Contract struct {
	address common.Address
	tx *types.Transaction
	contract bind.BoundContract
}

// DeployFn emulate store deploy function
type DeployFn func(...interface{})(common.Address, *types.Transaction, bind.BoundContract, error)

// LoadFn load contract
type LoadFn func(...interface{})(bind.BoundContract,error)


// ErrInstanceConversion failed type conversion
var (
	ErrInstanceConversion = errors.New("Can't convert type of smart contract")
	ErrLoadSc = errors.New("Invalid or imposible to load contract: %s \nError : %s")
)

// NewContract create a contract
func NewContract(
	address common.Address,
	tx *types.Transaction,
	contract bind.BoundContract,
) *Contract {
	return &Contract{address, tx, contract}
}

// JSON representation on contract
func (c *Contract) JSON() []byte {
	b, err := json.Marshal(helpers.Map{
		"address": c.address,
		"tx": c.tx,
	})

	if err != nil {
		return nil
	}

	return b
}
