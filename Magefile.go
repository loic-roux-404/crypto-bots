//+build mage

package main

import (
	"os"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	packageName = "github.com/loic-roux-404/crypto-bots"
)

var (
	ports = []string{"4205"}
	cmds  = []string{"sniper"}
	goexe = "go"
)

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")
}

type Build mg.Namespace

// Runs go mod download and then installs the binary.
func (Build) Api() error {
	/* if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
    return sh.Run("go", "install", "./...")
    */
    return nil
}

// Runs go mod download and then installs the binary.
func (Build) Web() error {
    return nil
}

// Runs go mod download and then installs the binary.
func (Build) Cmds() error {
    return nil
}

type Test mg.Namespace

// Runs go mod download and then installs the binary.
func (Test) Api() error {
    return nil
}

// Runs go mod download and then installs the binary.
func (Test) Web() error {
    return nil
}

// Runs go mod download and then installs the binary.
func (Test) Cmds() error {
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
    if (prerelease) {
        semverFlags = append(semverFlags, "prerelease")
    }

    if (noCi) {
        semverFlags = append(semverFlags, "no-ci")
    }

    finalFlags := strings.Split(strings.Join(semverFlags, " --"), " ")

	return sh.Run("semantic-release", finalFlags...)
}

type Deploy mg.Namespace

func (Deploy) Web() error {
	return nil
}

func (Deploy) Api() error {
	return nil
}

func (Deploy) Cmds() error {
	return nil
}

// Remove dev libraries and build/test artifacts
func Clean() error {
	return nil //clean.Clean()
}