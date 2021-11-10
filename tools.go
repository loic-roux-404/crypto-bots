//go:build tools
// +build tools

package tools

// Need this file to fix dependency resolve when installing go cmds
import (
	// etherum cmds import
	_ "github.com/ethereum/go-ethereum/cmd/abigen"
	_ "github.com/ethereum/go-ethereum/cmd/evm"
	_ "github.com/ethereum/go-ethereum/cmd/geth"
	// mage build tool import
	_ "github.com/magefile/mage"
)
