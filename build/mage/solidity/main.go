package solidity

import (
	"fmt"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
)

// Solidity mage namespace
type Solidity mg.Namespace

const solSelect = "solc-select"

// Install if missing
func (Solidity) install(version string) error {
	if helpers.CommandExist(solSelect) {
		return nil
	}

	if helpers.CommandExist("python3") || helpers.CommandExist("python") {
		return fmt.Errorf("Missing python3, install it to use solidity")
	}

	sh.Run("pip3", "install", solSelect)
	sh.Run(solSelect, "install", version)
	sh.Run(solSelect, "use", solSelect, version)

	return nil
}

// Compile smart contract
func (s Solidity) Compile(mockLoc string, mockName string, mockDest string) error {
	// TODO use a yaml config to set version
	err := s.install("0.5.16"); if err != nil {
		return err
	}

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

// Deploy a smart contract
func (s Solidity) Deploy(src string) error {
	return nil
}
