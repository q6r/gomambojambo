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
	"log"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

// generateAESDecryptAST will generate the decryption functin AST as a
// ast.CallExpr represent an anonymous function.
// eg : (func(s string) { ...... return plaintext })(encrypted_str)
func generateAESDecryptAST(key, nonce string) (*ast.CallExpr, error) {
	src := fmt.Sprintf(`package main
const a = (func(s string) string { k, _ := hex.DecodeString("%s"); ct, _ := hex.DecodeString(s); n, _ := hex.DecodeString("%s"); b, _ := aes.NewCipher(k); g, _ := cipher.NewGCM(b); pt, _ := g.Open(nil, n, ct, nil); return string(pt) })()`, key, nonce)

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		return nil, err
	}
	// TODO : meh
	decryptionCallExpr := f.Decls[0].(*ast.GenDecl).Specs[0].(*ast.ValueSpec).Values[0].(*ast.CallExpr)
	return decryptionCallExpr, nil
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
			newDecryptStrings(pkg.Name, fileast, key, nonce)
		}
	}

	// add imports to all pkgs packages
	for _, pkg := range pkgs {
		for _, fileast := range pkg.Files {
			astutil.AddImport(fset, fileast, "crypto/aes")
			astutil.AddImport(fset, fileast, "crypto/cipher")
			astutil.AddImport(fset, fileast, "encoding/hex")
		}
	}
	return nil
}

// newDecryptStrings will find all strings and wrap them in the anonymous
// function that will decrypt the string
func newDecryptStrings(pkgName string, fileAst *ast.File, key, nonce string) {

	astutil.Apply(fileAst, func(cr *astutil.Cursor) bool {
		cn, ok := cr.Node().(*ast.BasicLit)
		if !ok {
			return true
		}

		// We see an encrypted string, we wrap it in decryptionCallExpr
		if cn.Kind != token.STRING {
			return true
		}

		_, parentAssign := cr.Parent().(*ast.AssignStmt)
		_, parentIdent := cr.Parent().(*ast.ValueSpec)
		_, parentCallExpr := cr.Parent().(*ast.CallExpr)

		if parentAssign || parentIdent || parentCallExpr {
			// Insert the decryption function and required imports
			aesDecAST, err := generateAESDecryptAST(key, nonce)
			if err != nil {
				panic(err)
			}
			aesDecAST.Args = []ast.Expr{cn}
			cr.Replace(aesDecAST)

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

		_, parentAssign := cr.Parent().(*ast.AssignStmt)
		_, parentIdent := cr.Parent().(*ast.ValueSpec)
		_, parentCallExpr := cr.Parent().(*ast.CallExpr)

		if parentAssign || parentIdent || parentCallExpr {
			log.Printf("Enc:Assign : %#v, Current : %#v Parent : %#v\n", cn, cr.Node(), cr.Parent())
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

	// Never use more than 2^32 random nonces with a given key because of
	// the risk of a repeat.
	nonce, _ := hex.DecodeString(nonceHex)
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	return fmt.Sprintf("%x", ciphertext)
}
