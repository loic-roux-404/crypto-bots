package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
)

// CommandExist in PATH
func CommandExist(cmd string) bool {

	_, err := exec.LookPath(cmd)

	return err == nil
}

// GetCurrDir absolute path
func GetCurrDir() string {
	ex, err := os.Executable()

	if err != nil {
		panic(err)
	}

	return filepath.Dir(ex)
}