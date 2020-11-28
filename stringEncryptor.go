package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"

	log "github.com/sirupsen/logrus"

	"golang.org/x/tools/go/ast/astutil"
)

// generateAESDecryptAST will generate a decryption function as funcDecl for AES
// with key and nonce
func generateAESDecryptAST(key, nonce string) (*ast.FuncDecl, error) {
	src := fmt.Sprintf(`
	package main
	func AES_DECRYPT(s string) string {
		key, _ := hex.DecodeString("%s")
		ciphertext, _ := hex.DecodeString(s)
		nonce, _ := hex.DecodeString("%s")
		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err.Error())
		}

		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			panic(err.Error())
		}

		plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			panic(err.Error())
		}
		return string(plaintext)
	}`, key, nonce)

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

	return funcDecl, nil
}

// encryptString will encrypt all strings with key,nonce provided... insert the
// decryption function, and wrappers around encrypted strings
func encryptStrings(fset *token.FileSet, pkgs map[string]*ast.Package, key, nonce string) error {

	if len(key) != 64 {
		return errors.New("encryption key invalid length")
	}

	if len(nonce) != 24 {
		return errors.New("encryption nonce invalid length")
	}

	// encrypt strings to AES
	for _, pkg := range pkgs {
		for _, fileast := range pkg.Files {
			newEncryptStrings(pkg.Name, fileast, key, nonce)
		}
	}

	// wrap them with AES_DECRYPT call
	for _, pkg := range pkgs {
		for _, fileast := range pkg.Files {
			newDecryptStrings(pkg.Name, fileast)
		}
	}

	// Insert the decryption function and required imports
	aesDecAST, err := generateAESDecryptAST(key, nonce)
	if err != nil {
		return err
	}

	insertedFunction := false
	for _, pkg := range pkgs {
		for _, fileast := range pkg.Files {

			// Find main and insert decryption function before
			astutil.Apply(fileast, func(cr *astutil.Cursor) bool {
				// no need to search for main we already
				// inserted the decryption function
				if insertedFunction {
					return false
				}
				cn, ok := cr.Node().(*ast.FuncDecl)
				if !ok {
					return true
				}
				if cn.Name.String() != "main" {
					return true
				}
				cr.InsertBefore(aesDecAST)

				// Add required imports for aes decryption function
				astutil.AddImport(fset, fileast, "crypto/aes")
				astutil.AddImport(fset, fileast, "crypto/cipher")
				astutil.AddImport(fset, fileast, "encoding/hex")
				insertedFunction = true
				return false
			}, nil)
		}
	}

	return nil
}

// newDecryptStrings will insert a wrapper around encrypted strings to
// call the decryption function
func newDecryptStrings(pkgName string, fileAst *ast.File) {

	astutil.Apply(fileAst, func(cr *astutil.Cursor) bool {

		// find baselits
		// A BasicLit node represents a literal of basic type.
		cn, ok := cr.Node().(*ast.BasicLit)
		if !ok {
			return true
		}

		if cn.Kind != token.STRING {
			return true
		}

		assignv, parentAssignOk := cr.Parent().(*ast.AssignStmt)
		identv, parentIdentOk := cr.Parent().(*ast.ValueSpec)
		callv, parentCallExprOk := cr.Parent().(*ast.CallExpr)

		// If a basic lit "string is found", we search the ast
		// for its parent, convert it to GenDecl and determine
		// the token type and we set isConst to that, so we
		// avoid doing something like :
		// const MyConst = AES_DECRYPT("....")
		// because calls not allowed in const
		isConst := false
		astutil.Apply(fileAst, func(cr *astutil.Cursor) bool {
			if (cr.Node() == identv) || (cr.Node() == assignv) || (cr.Node() == callv) {
				gendec, ok := cr.Parent().(*ast.GenDecl)
				if !ok {
					return true
				}
				isConst = gendec.Tok == token.CONST
				return false
			}
			return true
		}, nil)

		if (parentAssignOk || parentIdentOk || parentCallExprOk) && !isConst {

			cr.Replace(&ast.CallExpr{
				Fun:  ast.NewIdent("AES_DECRYPT"),
				Args: []ast.Expr{cn},
			})
		}

		return true
	}, nil)
}

// newEncryptStrings will encrypt all strings with key,nonce
func newEncryptStrings(pkgName string, fileAst *ast.File, key, nonce string) {
	astutil.Apply(fileAst, func(cr *astutil.Cursor) bool {
		cn, ok := cr.Node().(*ast.BasicLit)
		if !ok {
			return true
		}

		if cn.Kind != token.STRING {
			return true
		}

		assignv, parentAssign := cr.Parent().(*ast.AssignStmt)
		identv, parentIdent := cr.Parent().(*ast.ValueSpec)
		callv, parentCallExpr := cr.Parent().(*ast.CallExpr)

		// If a basic lit "string is found", we search the ast
		// for its parent, convert it to GenDecl and determine
		// the token type and we set isConst to that, so we
		// avoid doing something like :
		// const MyConst = AES_DECRYPT("....")
		// because calls not allowed in const
		isConst := false
		astutil.Apply(fileAst, func(cr *astutil.Cursor) bool {
			if (cr.Node() == identv) || (cr.Node() == assignv) || (cr.Node() == callv) {
				gendec, ok := cr.Parent().(*ast.GenDecl)
				if !ok {
					return true
				}
				isConst = gendec.Tok == token.CONST
				return false
			}
			return true
		}, nil)

		if (parentAssign || parentIdent || parentCallExpr) && !isConst {
			if *verbose {
				log.Printf("Enc:Assign : %#v, Current : %#v Parent : %#v\n", cn, cr.Node(), cr.Parent())
			}
			valInterpreted, err := strconv.Unquote(cn.Value)
			if err != nil {
				panic(err)
			}
			cr.Replace(&ast.BasicLit{
				Value: fmt.Sprintf("\"%s\"", encryptString(valInterpreted, key, nonce)),
				Kind:  token.STRING,
			})
		}

		return true
	}, nil)
}

func encryptString(plaintext string, keyHex, nonceHex string) string {
	key, _ := hex.DecodeString(keyHex)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce, _ := hex.DecodeString(nonceHex)
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	return fmt.Sprintf("%x", ciphertext)
}
