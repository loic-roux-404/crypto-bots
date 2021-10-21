package cmd

import (
	"fmt"
	"github.com/loic-roux-404/crypto-bots/tools/system"
	"strings"

	"github.com/magefile/mage/sh"
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
	return system.CommandExist(pkgToCmd(pkg))
}

// BinInstall install bins
func BinInstall(tools []string) error {
	for _, pkg := range tools {

		if pkgCommandExist(pkg) {
			fmt.Printf("%s Already installed", pkg)
			return nil
		}

		err := sh.Run("go", "install", pkg)

		if err != nil {
			return err
		}
	}

	return nil
}