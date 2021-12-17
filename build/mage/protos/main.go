package protos

import (
	"log"
	"os"

	"github.com/magefile/mage/sh"
)

func BufAll(name string) (err error) {
	os.Chdir("api")
	// Empty for all proto files
	var args []string = []string{}

	if name == "" {
		log.Print("Info: Building all")
	} else {
		args = []string{name}
		log.Printf("Info: Building %s", name)
	}

	err = bufCmd("build", args); if err != nil {
		return err
	}

	err = bufCmd("generate", args); if err != nil {
		return err
	}

	return nil
}

func bufCmd(subCmd string, args []string) error {
	return sh.RunV("buf", append([]string{subCmd}, args...)...)
}
