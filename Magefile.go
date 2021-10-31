//go:build mage
// +build mage

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"

	"github.com/loic-roux-404/crypto-bots/build/mage/cmd"
	// mage:import
	_ "github.com/loic-roux-404/crypto-bots/build/mage/release"
	// mage:import
	"github.com/loic-roux-404/crypto-bots/build/mage/solidity"
	"github.com/loic-roux-404/crypto-bots/internal/helpers"
)

const (
	packageName = "github.com/loic-roux-404/crypto-bots"
)

var (
	ports      = []string{"4205"}
	cmds       = []string{"sniper"}
	goexe      = "go"
	currentDir = helpers.GetCurrDir()
	binDir, _ = filepath.Abs(filepath.Join(".", "bin"))
)

var toolsCmds = []string{
	"github.com/ethereum/go-ethereum/cmd/evm",
	"github.com/ethereum/go-ethereum/cmd/geth",
	"github.com/ethereum/go-ethereum/cmd/abigen",
}

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")
	// Etherum need a c compiler
	// Verify if a clang / gcc exist in your PATH
	os.Setenv("CGO_ENABLED", "1")

	err := cmd.BinInstall(toolsCmds)

	if err != nil {
		panic(err)
	}
}

var Default = All

type Build mg.Namespace

// All Build all modules
func All() error {
	b := new(Build)

	err := b.ScPancake(); if err != nil {
		return err
	}

	err = b.Cmds(""); if err != nil {
		return err
	}

	err = b.Api(); if err != nil {
		return err
	}

	err = b.Web(); if err != nil {
		return err
	}

	return nil
}

// Api build protobuf files
func (Build) Api() error {
	return nil
}

// Web interface build
func (Build) Web() error {
	return nil
}

var cmdCompiler = cmd.NewCompiler(goexe, cmds, "cmd/bot", binDir)

// Cmds build all CLI
func (Build) Cmds(name string) error {
	return cmdCompiler.GoexeCmd(name)
}

// CmdsRun dev command
func (Build) CmdsRun(name string) error {
	return cmdCompiler.GoexeRun(name)
}

// ScPancake compile smart contract and generate library
func (Build) ScPancake() error {
	s := new(solidity.Solidity)
	if err := s.Compile(mockLoc, mockName, mockDest); err != nil {
		log.Printf("Warn: %v", err)
	}

	if err := s.PackageByNet(scByNetSet, "gencontracts"); err != nil {
		return err
	}

	return nil
}

type Test mg.Namespace

var (
	mockLoc  = filepath.Join(".", "tests", "mocks")
	mockDest = filepath.Join(mockLoc, "data")
	mockName = "glitch"
	scByNetSet = helpers.Map{"erc20": filepath.Join(mockDest, "PancakePair")}
)

// Runs go mod download and then installs the binary.
func (Test) Api() error {
	return nil
}

// Web Test
// Runs go mod download and then installs the binary.
// e2e test TODO (k6) ?
func (Test) Web() error {
	return nil
}

// Runs go mod download and then installs the binary.
func (t Test) Cmds() error {

	return nil
}

type Deploy mg.Namespace

func (Deploy) ScPancake() error {
	return nil
}

// Web ui deploy
func (Deploy) Web() error {
	return nil
}

// Api grpc endpoint
func (Deploy) Api() error {
	return nil
}

// Command line deploy
func (Deploy) Cmds() error {
	return nil
}

// Clean Remove dev libraries and build/test artifacts
func Clean() error {
	return nil //clean.Clean()
}
