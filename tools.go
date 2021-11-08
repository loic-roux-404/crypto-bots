package tools

// Import cmd builded packages
import (
	// etherum cmds import
	_ "github.com/ethereum/go-ethereum/cmd/abigen"
	_ "github.com/ethereum/go-ethereum/cmd/evm"
	_ "github.com/ethereum/go-ethereum/cmd/geth"
	// build tool import
	_ "github.com/magefile/mage"
)
