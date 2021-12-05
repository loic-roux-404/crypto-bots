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

func (f Fmt) All() (err error) {
	err = f.Go(); if err != nil {
		return err
	}

	err = f.ProtoCheck("grpc/proto", "buf lint"); if err != nil {
		return err
	}

	return nil
}

// GoFix bad formatted files
func (Fmt) GoFix() error {
	return sh.Run("gofmt", "-s", "-d", "-w", ".")
}

// Go go files format (useful in a CI)
func (Fmt) Go() error {
	out, err := helpers.RunAndGetStdout(sh.RunWithV, env, "gofmt", "-d", "-e", "-l", ".")

	return fmtErrExit(out, "Error: Syntax in file : %s", err)
}

func (Fmt) ProtoCheck(dir, bin string) error {
	out, err := helpers.RunAndGetStdout(sh.RunWithV, env, "buf", "lint", dir)

	return fmtErrExit(out, "Error: Syntax in proto file : %s", err)
}

func fmtErrExit(out, msgFmt string, err error) error {
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
