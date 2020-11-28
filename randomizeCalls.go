package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"math/rand"
	"strings"
	"unicode"

	"golang.org/x/tools/go/ast/astutil"
)

func randomizeCalls(fset *token.FileSet, pkgs map[string]*ast.Package) map[string]string {

	funcChangeHistory := make(map[string]string)

	// randomize the top level functions
	for _, pkg := range pkgs {
		for _, fileast := range pkg.Files {
			newRandomizeTop(pkg.Name, fileast, funcChangeHistory)
		}
	}

	// randomize all top level function bodies
	for _, pkg := range pkgs {
		for _, fileast := range pkg.Files {
			newRandomizeInner(pkg.Name, fileast, funcChangeHistory)
		}
	}

	return funcChangeHistory
}

func newRandomizeInner(pkgName string, fileAst *ast.File, changeHistory map[string]string) {
	astutil.Apply(fileAst, func(cr *astutil.Cursor) bool {
		callExpr, ok := cr.Node().(*ast.CallExpr)
		if !ok {
			return true
		}

		switch fun := callExpr.Fun.(type) {
		case *ast.SelectorExpr:
			ident, ok := fun.X.(*ast.Ident)
			if ok {
				if rname, ok := changeHistory[fmt.Sprintf("%s.%s", ident.Name, fun.Sel.Name)]; ok {

					fun.Sel = ast.NewIdent(rname)
				}
			}
		case *ast.Ident:
			if rname, ok := changeHistory[fmt.Sprintf("%s.%s", pkgName, fun.Name)]; ok {
				fun.Name = rname
			}
		}

		return true
	}, nil)
}

func newRandomizeTop(pkgName string, fileAst *ast.File, changeHistory map[string]string) {
	astutil.Apply(fileAst, func(cr *astutil.Cursor) bool {
		funcDecl, ok := cr.Node().(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Ignore functions with recievers
		if funcDecl.Recv != nil {
			return true
		}

		// Ignore main
		if funcDecl.Name.String() == "main" && pkgName == "main" {
			return true
		}

		// Ignore init function
		if funcDecl.Name.String() == "init" {
			return true
		}

		outname := fmt.Sprintf("%s.%s", pkgName, funcDecl.Name.String())

		// if it already exists
		if randomName, ok := changeHistory[outname]; ok {
			funcDecl.Name = ast.NewIdent(randomName)
		} else {
			randomName := randStringRunes(32)
			if isExportedFunction(string(funcDecl.Name.String())) {
				randomName = strings.Title(randomName)
			}
			changeHistory[outname] = randomName
			funcDecl.Name = ast.NewIdent(randomName)
		}

		return true
	}, nil)
}

func extractRecvTypeFromFuncDecl(funcDecl *ast.FuncDecl) (string, error) {
	typeIdent := ""
	if funcDecl.Recv == nil {
		return "", errors.New("no recieve on function")
	}

	for _, field := range funcDecl.Recv.List {
		astutil.Apply(field.Type, func(cr *astutil.Cursor) bool {
			ident, ok := cr.Node().(*ast.Ident)
			if !ok {
				return true
			}
			typeIdent = ident.Name
			return false
		}, nil)
	}

	if typeIdent == "" {
		return "", errors.New("unable to extract type ident")
	}
	return typeIdent, nil
}

func isExportedFunction(funcName string) bool {
	if len(funcName) == 0 {
		panic("this should not happen")
		return false
	}

	return unicode.IsUpper(rune(funcName[0]))
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
