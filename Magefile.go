//go:build mage
// +build mage

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"

	"github.com/loic-roux-404/crypto-bots/build/mage/cmd"
	"github.com/loic-roux-404/crypto-bots/build/mage/goutils"
	"github.com/loic-roux-404/crypto-bots/build/mage/protos"

	// mage:import
	ci "github.com/loic-roux-404/crypto-bots/build/mage/ci"
	// mage:import
	"github.com/loic-roux-404/crypto-bots/build/mage/solidity"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
)

const (
	packageName = "github.com/loic-roux-404/crypto-bots"
)

var (
	ports       = []string{"4205"}
	cmds        = []string{"cryptos-bots"}
	goexe       = "go"
	currentDir  = helpers.GetCurrDir()
	binDir, _   = filepath.Abs(filepath.Join(".", "bin"))
	env         = map[string]string{}
	solcVersion = "0.5.16"
)

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

	imports, err := goutils.ImportsToSlice(".", "tools.go", "tools")

	if err != nil {
		panic(err)
	}

	err = cmd.BinInstall(imports)
	if err != nil {
		panic(err)
	}
}

var Default = All

type Build mg.Namespace

// All Build all modules
func All() error {
	b := new(Build)

	err := b.ScPancake()
	if err != nil {
		return err
	}

	err = b.Cmds("")
	if err != nil {
		return err
	}

	err = b.Api()
	if err != nil {
		return err
	}

	return b.App()
}

// Api build protobuf files
func (Build) Api() error {
	return protos.BufAll("")
}

// App interface build
func (Build) App() error {
	return nil
}

// Commands
var (
	cmdCompiler = cmd.NewCompiler(goexe, cmds, "cmd", binDir)
)

// Contracts
var (
	gencontractsDir = "gencontracts"
	mockLoc         = filepath.Join(".", "tests", "mocks")
	mockDest        = filepath.Join(mockLoc, "data")
	mockName        = "glitch"
	scByNetSet      = helpers.SimpleMap{"bep20": filepath.Join(mockDest, "PancakePair")}
	// unit tests options
	unitTimeout = "30s"
	// Create test runnner modules
	testRunner = ci.NewRunner(goexe, packageName, env, map[string]string{
		"timeout": unitTimeout,
	})
)

// Cmds build all CLI
func (Build) Cmds(name string) error {
	return cmdCompiler.GoexeCmd(name, env)
}

// CmdsRun dev command
func (Build) CmdsRun(name string) error {
	return cmdCompiler.GoexeRun(name, env)
}

// ScPancake compile smart contract and generate library
func (Build) ScPancake() error {
	s, err := solidity.NewSolidity(solcVersion)
	if err != nil {
		return err
	}

	err = s.Compile(mockLoc, mockName, mockDest)
	if err != nil {
		log.Printf("Warn: %v", err)
	}

	return s.PackageByNet(scByNetSet, gencontractsDir)
}

type Test mg.Namespace

// TODO set tests folders
func (t Test) All() (err error) {
	err = t.Lib("")
	if err != nil {
		return err
	}

	err = t.App()
	if err != nil {
		return err
	}

	return t.Api()
}

// Api end to end and load testing
func (Test) Api() error {
	return nil
}

// App Test
// Forward to all node js tests
// cypress / jest / BDD
func (Test) App() error {
	return nil
}

// Lib Tests, all module or single one.
// Libraries are used in cmd, pkg and api
func (Test) Lib(pkg string) error {
	return testRunner.Pkg(pkg)
}

type Deploy mg.Namespace

func (Deploy) ScPancake() error {
	return nil
}

// App ui deploy
func (Deploy) App() error {
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
