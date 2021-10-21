package solidity

import (
	"fmt"
	"path/filepath"

	"github.com/magefile/mage/sh"
)

// Compile smart contract
func Compile(mockLoc string, mockName string, mockDest string) error {
	finalMockName := fmt.Sprintf("%s.sol", mockName)
	var (
		pcv2Contract = filepath.Join(mockLoc, finalMockName)
		pcv2Bin = filepath.Join(mockDest, fmt.Sprintf("%s.bin", mockName))
		pcv2Abi = filepath.Join(mockDest, fmt.Sprintf("%s.abi", mockName))
		mockDestArg = fmt.Sprintf("--output-dir=%s", mockDest)
	)


	if err := sh.Run("solc", "--abi", pcv2Contract, mockDestArg);
	err != nil {
		return err
	}

	if err := sh.Run("solc", "--bin", pcv2Contract, mockDestArg);
	err != nil {
		return err
	}

	return sh.Run(
		"abigen", 
		fmt.Sprint("--bin=", pcv2Bin), 
		fmt.Sprint("--abi=", pcv2Abi), 
		fmt.Sprint("--pkg=", mockDest),
		fmt.Sprint("--out=", fmt.Sprintf("%s.go", mockName)),
	)
}