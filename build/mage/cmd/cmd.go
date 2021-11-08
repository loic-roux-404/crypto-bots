package cmd

import (
	"log"
	"strings"

	"github.com/magefile/mage/sh"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
)

// ToFlags transform
func ToFlags(flagList []string) []string {
	return strings.Split(strings.Join(flagList, " --"), " ")
}

// ToFlagsStr transform
func ToFlagsStr(flagList []string) string {
	return strings.Join(flagList, " --")
}

func pkgToCmd(pkg string) string {
	parts := strings.Split(pkg, "/")

	return parts[len(parts)-1]
}

// PkgCommandExist in PATH
func pkgCommandExist(pkg string) bool {
	return helpers.CommandExist(pkgToCmd(pkg))
}

// BinInstall install bins
func BinInstall(tools []string) error {
	for _, pkg := range tools {

		if pkgCommandExist(pkg) {
			// Verbose mode ?
			continue
		}

		log.Printf("Info: Installing %s\n", pkg)

		err := sh.Run("go", "install", pkg)

		if err != nil {
			return err
		}
	}

	return nil
}

// ToLocalCmd to local command
func ToLocalCmd(prefix string, cmd string) string {
	return strings.Join([]string{prefix, cmd, "main.go"}, "/")
}
