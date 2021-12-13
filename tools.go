// +build tools

package tools

// Need this file to fix dependency resolve when installing go cmds
// This var is processed by mage
import (
	// etherum cmds import
	_ "github.com/ethereum/go-ethereum/cmd/abigen"
	_ "github.com/ethereum/go-ethereum/cmd/evm"
	_ "github.com/ethereum/go-ethereum/cmd/geth"
	// mage build tool import
	_ "github.com/magefile/mage"
	// Cobra
	_ "github.com/spf13/cobra/cobra"
	// Proto imports
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking"
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-lint"
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
