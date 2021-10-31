package cmd

import (
	"fmt"

	"github.com/magefile/mage/sh"
	"github.com/thoas/go-funk"
)

const (
	// BUILD mode
	BUILD = "build"
	// RUN mode
	RUN = "run"
	// TEST mode
	TEST = "test"
)

// Compiler i
type Compiler struct{
	cmds []string
	goexe string
	src string
	dest string
	mode string
}

// NewCompiler contruct commile module
func NewCompiler(goexe string, cmds []string, src string, dest string) *Compiler {
	mode := BUILD

	if (len(goexe)) <= 0 {
		goexe = "go"
	}

	return &Compiler{
		cmds,
		goexe,
		src,
		dest,
		mode,
	}
}

// GoexeCmd Run a go compiler command
// Default is build
func (c *Compiler) GoexeCmd(name string) error {
	if len(name) > 0 {
		c.cmds = funk.Filter(c.cmds, func(x string) bool {
			return x == name
		}).([]string)
	}

	for _, cmd := range c.cmds {
		finalc := ToLocalCmd(c.src, cmd)
		// Default to verbose in run (dev) mode
		dest := []string{c.mode}

		if c.mode == BUILD {
			c.dest = fmt.Sprintf("%s/%s", c.dest, cmd)
			fmt.Printf("Building %s in %s...\n", finalc, c.dest)
			dest = append(dest, "-o", c.dest)
		}

		dest = append(dest, finalc)
		err := sh.Run(c.goexe, dest...)

		if err != nil {
			return err
		}
	}

	return nil
}

// GoexeRun go run command
func (c *Compiler) GoexeRun(name string) error {
	c.mode = RUN
	err := c.GoexeCmd(name)
	c.mode = BUILD

	return err
}
