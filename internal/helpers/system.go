package helpers

import (
	"errors"
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

// Exists file
func Exists(name string) (bool, error) {
    _, err := os.Stat(name)
    if err == nil {
        return true, nil
    }
    if errors.Is(err, os.ErrNotExist) {
        return false, nil
    }
    return false, err
}
