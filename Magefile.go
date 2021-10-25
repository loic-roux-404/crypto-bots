//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/loic-roux-404/crypto-bots/build/mage/cmd"
	// mage:import
	"github.com/loic-roux-404/crypto-bots/build/mage/solidity"
	"github.com/loic-roux-404/crypto-bots/internal/system"
	"github.com/thoas/go-funk"
)

const (
	packageName = "github.com/loic-roux-404/crypto-bots"
)

var (
	ports      = []string{"4205"}
	cmds       = []string{"sniper"}
	goexe      = "go"
	currentDir = system.GetCurrDir()
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

// Runs go mod download and then installs the binary.
func All() error {
	b := new(Build)

	err := b.Cmds("")
	if err != nil {
		return err
	}

	err = b.Api()
	if err != nil {
		return err
	}

	err = b.Web()
	if err != nil {
		return err
	}

	return nil
}

// Runs go mod download and then installs the binary.
func (Build) Api() error {
	return nil
}

// Runs go mod download and then installs the binary.
func (Build) Web() error {
	return nil
}

const (
	BUILD = "build"
	RUN   = "run"
)

var mode = BUILD

// Runs go mod download and then installs the binary.
func (Build) Cmds(name string) error {
	if len(name) > 0 {
		cmds = funk.Filter(cmds, func(x string) bool {
			return x == name
		}).([]string)
	}

	for _, c := range cmds {
		finalc := cmd.ToLocalCmd("cmd/bot", c)
		// Default to verbose in run (dev) mode
		dest := "-v"

		if mode == BUILD {
			dest = filepath.Join(".", "bin", c)
			fmt.Printf("Building %s in %s...", finalc, dest)
			dest = fmt.Sprintf("-o %s", dest)
		}

		err := sh.Run(goexe, mode, dest, finalc)
		if err != nil {
			return err
		}
	}

	return nil
}

// CmdsRun dev command
func (b Build) CmdsRun(name string) error {
	mode = RUN
	err := b.Cmds(name)
	mode = BUILD

	return err
}

type Test mg.Namespace

var (
	mockLoc  = filepath.Join(".", "tests", "mocks")
	mockDest = filepath.Join(mockLoc, "data")
	mockName = "glitch"
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
	s := new(solidity.Solidity)
	if err := s.Compile(mockLoc, mockName, mockDest); err != nil {
		return err
	}

	return nil
}

type Release mg.Namespace

const (
	semverExe = "semantic-release"
)

var semverFlags = []string{
	"allow-initial-development-versions",
	"download-plugins",
}

// Release semantic release
func (Release) SemRelease(prerelease bool, noCi bool) error {
	if prerelease {
		semverFlags = append(semverFlags, "prerelease")
	}

	if noCi {
		semverFlags = append(semverFlags, "no-ci")
	}

	finalFlags := strings.Split(strings.Join(semverFlags, " --"), " ")

	return sh.Run("semantic-release", finalFlags...)
}

type Deploy mg.Namespace

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
