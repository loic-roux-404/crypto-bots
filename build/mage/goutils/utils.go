package goutils

import (
	// "go/ast"
	// "go/importer"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	// "go/types"
	"os"
	"path/filepath"
)

var (
	ErrExpectedDir = errors.New("entry point should be a path to a directory of a go package")
)

// TODO unit test
func ImportsToSlice(pkgDir, file, pkgName string) ([]string, error) {
	pkgs, err := PkgsFromFile(pkgDir, file)

	if err != nil {
		return nil, err
	}

	imports := []string{}
	if len(pkgs) <= 0 {
		return imports, nil
	}

	for _, f := range pkgs[pkgName].Files {
		for _, pkgImport := range f.Imports {
			imports = append(imports, strings.Replace(pkgImport.Path.Value, "\"", "", -1))
		}
	}

	return imports, nil
}

func PkgsFromFile(relDir, file string) (map[string]*ast.Package, error) {
	fp, err := getPkgAbsDir(relDir)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(
		fset,
		fp,
		func(info os.FileInfo) bool {
			return info.Name() == file
		},
		parser.ImportsOnly,
	)

	if err != nil {
		return nil, err
	}

	return pkgs, nil
}

func getPkgAbsDir(relDir string) (string, error) {
	fp, err := filepath.Abs(relDir)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(fp)

	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", ErrExpectedDir
	}

	return fp, nil
}
