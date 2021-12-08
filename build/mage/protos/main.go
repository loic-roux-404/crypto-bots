package protos

import (
	"log"
	"os"

	"github.com/magefile/mage/sh"
)

func BufAll(name string) (err error) {
	os.Chdir("api")
	// Empty for all
	var arg []string = []string{}

	if name == "" {
		log.Print("Info: Building all")
	} else {
		arg = []string{name}
		log.Printf("Info: Building %s", name)
	}

	err = bufCmd("build", arg); if err != nil {
		return err
	}

	err = bufCmd("generate", arg); if err != nil {
		return err
	}

	return nil
}

func bufCmd(subCmd string, args []string) error {
	return sh.RunV("buf", append([]string{subCmd}, args...)...)
}
