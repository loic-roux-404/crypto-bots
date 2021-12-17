package kecacc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ErrAddressInvalid invalid
var (
	ErrAddressInvalid = errors.New("invalid public address")
	ErrScInvalid  = errors.New("invalid smart contract address")
)

// ValidateAddress destination
func ValidateAddress(address interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := address.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// ValidateAccAddress destination
func ValidateAccAddress(acc accounts.Account) bool {
	return ValidateAddress(acc.Address)
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
	return err == nil
}

// IsErrAddress validate address but return an error if invalid
func IsErrAddress(address common.Address) error {
	acc := accounts.Account{
		Address: address,
	}
	isValidAd := ValidateAccAddress(acc)

	if isValidAd {
		return nil
	}

	return fmt.Errorf("%s : %s", ErrAddressInvalid, address)
}

// IsErrStrAddress is this string address ok
func IsErrStrAddress(address string) error {
	finalAddress := common.HexToAddress(address)

	return IsErrAddress(finalAddress)
}

// IsZeroAddress validate if it's a 0 address
func IsZeroAddress(iaddress interface{}) bool {
	var address common.Address
	switch v := iaddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}

// SigRSV signatures R S V returned as arrays
func SigRSV(isig interface{}) ([32]byte, [32]byte, uint8) {
	var sig []byte
	switch v := isig.(type) {
	case []byte:
		sig = v
	case string:
		sig, _ = hexutil.Decode(v)
	}

	sigstr := common.Bytes2Hex(sig)
	rS := sigstr[0:64]
	sS := sigstr[64:128]
	R := [32]byte{}
	S := [32]byte{}
	copy(R[:], common.FromHex(rS))
	copy(S[:], common.FromHex(sS))
	vStr := sigstr[128:130]
	vI, _ := strconv.Atoi(vStr)
	V := uint8(vI + 27)

	return R, S, V
}
