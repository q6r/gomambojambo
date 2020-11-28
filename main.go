package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

var srcPath = flag.String("srcpath", "", "the path to the src directory")
var writeChanges = flag.Bool("writechanges", false, "write changes to files")
var obfCalls = flag.Bool("calls", false, "enable randomization of calls and functions")
var obfLoops = flag.Bool("loops", false, "obfuscate loops by converting to gotos")
var obfStrings = flag.Bool("strings", false, "obfuscate strings by encryption")
var obfStringKey = flag.String("stringsKey", "0101010101010101010101010101010101010101010101010101010101010101", "the key for encrypting strings (64 length)")
var obfStringNonce = flag.String("stringNonce", "010101010101010101010101", "the nonce for encrypting strings (24 length)")
var verbose = flag.Bool("verbose", false, "be verbose")
var deadcode = flag.Bool("deadcode", false, "add some deadcode")

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func init() {
	rand.Seed(time.Now().UnixNano())
	log.SetOutput(os.Stdout)
}

func parseDirRecursive(fset *token.FileSet, path string, filter func(string) bool, mode parser.Mode) (pkgs map[string]*ast.Package, first error) {

	list := []string{}
	err := filepath.Walk(path, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		list = append(list, fpath)
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	pkgs = make(map[string]*ast.Package)
	for _, filename := range list {
		if strings.HasSuffix(filename, ".go") && (filter == nil || filter(filename)) {
			if *verbose {
				fmt.Printf("Parsing %#v\n", filename)
			}
			if src, err := parser.ParseFile(fset, filename, nil, mode); err == nil {
				name := src.Name.Name
				pkg, found := pkgs[name]
				if !found {
					pkg = &ast.Package{
						Name:  name,
						Files: make(map[string]*ast.File),
					}
					pkgs[name] = pkg
				}
				pkg.Files[filename] = src
			} else if first == nil {
				first = err
			}
		}
	}

	return
}

func main() {
	flag.Parse()

	if srcPath == nil {
		panic("provide --srcpath")
	}

	fset := token.NewFileSet()
	pkgs, err := parseDirRecursive(fset, *srcPath, func(d string) bool {
		return true
	}, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	if *deadcode {
		if err := injectDeadcode(fset, pkgs); err != nil {
			panic(err)
		}
	}

	if *obfStrings {
		if err := encryptStrings(fset, pkgs, *obfStringKey, *obfStringNonce); err != nil {
			panic(err)
		}
	}

	if *obfLoops {
		forToTagLoops(fset, pkgs)
	}

	if *obfCalls {
		funcChangeHistory := randomizeCalls(fset, pkgs)
		if *verbose {
			log.Printf("Functions randomized : %v", funcChangeHistory)
		}
	}

	// Show/Write changes
	for _, pkg := range pkgs {
		for file, fileast := range pkg.Files {
			buf := new(bytes.Buffer)
			if err := format.Node(buf, fset, fileast); err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", buf.Bytes())
			if *writeChanges {
				ioutil.WriteFile(file, buf.Bytes(), 0644)
			}
		}
	}
}
