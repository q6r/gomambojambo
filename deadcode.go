package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"

	log "github.com/sirupsen/logrus"

	"golang.org/x/tools/go/ast/astutil"
)

func generateDeadCode() ([]ast.Stmt, error) {
	xMin := 1
	xMax := 10
	zXXXinit := rand.Intn(xMax-xMin+1) + xMin
	sXXXinit := rand.Intn(xMax-xMin+1) + xMin
	iXXXinit := rand.Intn(xMax-xMin+1) + xMin

	// BUG : need to have it in anonymous function otherwwise
	// loop-obfuscator fails
	src := fmt.Sprintf(`
	package main
	func SOMEDEADCODE() {
		(func() {
			zXXX := int64(%d)	
			sXXX := float64(%d)
			for iXXX := %d; iXXX < 15; iXXX++ {
				for jXXX := iXXX; jXXX < 15; jXXX++ {
					for zXXX := jXXX; zXXX < 15; zXXX++ {
						sXXX = (float64(iXXX+ jXXX) * float64(zXXX)) / float64(iXXX);
					}
				}
			}
			if sXXX == float64(zXXX) {
				;
			}
		})()
	}`, zXXXinit, sXXXinit, iXXXinit)

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		return nil, err
	}

	// Extract function from Decls
	funcDecl, ok := f.Decls[0].(*ast.FuncDecl)
	if !ok {
		return nil, errors.New("failed to cast funcDecl")
	}

	return funcDecl.Body.List, nil
}

func injectDeadcode(fset *token.FileSet, pkgs map[string]*ast.Package) error {

	deadcodeAST, err := generateDeadCode()
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for _, fileast := range pkg.Files {
			astutil.Apply(fileast, func(cr *astutil.Cursor) bool {
				cn, ok := cr.Node().(*ast.FuncDecl)
				if !ok {
					return true
				}

				if *verbose {
					log.Printf("Adding deadcode to %s", cn.Name.Name)
				}
				cn.Body.List = append(deadcodeAST, cn.Body.List...)
				return true
			}, nil)
		}
	}

	return nil
}
