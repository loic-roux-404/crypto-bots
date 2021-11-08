package helpers

import (
	"bytes"
	"errors"
	"io"
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

// RunAndGetStdout in string format
func RunAndGetStdout(
	runFn func(env map[string]string, cmd string, args ...string) error,
	env map[string]string,
	cmd string,
	args ...string,
) (string, error) {
    old := os.Stdout // keep backup of the real stdout
    r, w, err := os.Pipe(); if err != nil {
		return "", err
	}

	os.Stdout = w

	err = runFn(env, cmd, args...); if err != nil {
		return "", err
	}

    outC := make(chan string)
    // copy the output in a separate goroutine so printing can't block indefinitely
    go func() {
        var buf bytes.Buffer
        io.Copy(&buf, r)
        outC <- buf.String()
    }()

    // back to normal state
    w.Close()
    os.Stdout = old // restoring the real stdout
	out := <-outC

	return out, nil
}
