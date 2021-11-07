package ci

import (
	"fmt"
	"os"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Fmt namespace
type Fmt mg.Namespace

var env = map[string]string{}

// Fix go files format
func (Fmt) Fix() error {
	return sh.Run("gofmt", "-s", "-d", "-w", ".")
}

// Check go files format (useful in a CI)
func (Fmt) Check() (err error) {
	out, err := helpers.RunAndGetStdout(sh.RunWithV, env, "gofmt", "-d", "-e", "-l", ".");
	hasFmtErr := len(out) > 0

	if !hasFmtErr && err == nil {
		return nil
	}

	if hasFmtErr {
		fmt.Printf("Error: Syntax in file : %s", out)
		os.Exit(1)
	}

	return err
}
