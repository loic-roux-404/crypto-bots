package ci

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// Runner test info
type Runner struct {
	goexe string
	packageName string
	env map[string]string
	args map[string]string
}

// NewRunner testing module
func NewRunner(goexe, packageName string, env, args map[string]string) *Runner {
	return &Runner{
		goexe,
		packageName,
		env,
		args,
	}
}

// Pkg test
// TODO transform map in slice
func (app *Runner) Pkg(pkg string) error {
	var fullPkg string

	if len(pkg) <= 0 {
		fullPkg = "./..."
	} else {
		fullPkg = fmt.Sprintf("%s/%s", app.packageName, pkg)
	}

	return sh.RunWithV(app.env, app.goexe, "test", "-timeout", "30s", fullPkg, "-v")
}
