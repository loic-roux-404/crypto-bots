//+build mage

package main

import (
	"os"
	"strings"

	"github.com/magefile/mage/sh"
)

const (
	packageName = "github.com/loic-roux-404/crypto-bots"
)

var goexe = "go"

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")
}

// Runs go mod download and then installs the binary.
func Sniper() error {
	/* if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
    return sh.Run("go", "install", "./...")
    */
    return nil
}

// TODO build deps
func Package() error {
    return nil
}

const (
    semverExe = "semantic-release"
)

var semverFlags = []string{
    "allow-initial-development-versions",
    "download-plugins",
}

// Release semantic release
func Release(prerelease bool, noCi bool) error {
    if (prerelease) {
        semverFlags = append(semverFlags, "prerelease")
    }

    if (noCi) {
        semverFlags = append(semverFlags, "no-ci")
    }

    finalFlags := strings.Split(strings.Join(semverFlags, " --"), " ")

	return sh.Run("semantic-release", finalFlags...)
}
