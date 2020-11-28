package main

import (
	"fmt"
	"go/ast"
	"go/token"

	log "github.com/sirupsen/logrus"

	"golang.org/x/tools/go/ast/astutil"
)

func forToTagLoops(fset *token.FileSet, pkgs map[string]*ast.Package) map[string]string {

	funcChangeHistory := make(map[string]string)

	// randomize all top level function bodies
	for _, pkg := range pkgs {
		for _, fileast := range pkg.Files {
			convertTopLevelFunctionsBodiesLoops(pkg.Name, fileast, funcChangeHistory)
		}
	}

	return funcChangeHistory
}

func convertTopLevelFunctionsBodiesLoops(pkgName string, fileAst *ast.File, changeHistory map[string]string) {
	astutil.Apply(fileAst, func(cr *astutil.Cursor) bool {
		stmtType, ok := cr.Node().(*ast.ForStmt)
		if !ok {
			return true
		}

		init := stmtType.Init
		cond := stmtType.Cond
		post := stmtType.Post
		body := stmtType.Body

		// We only support IfStmt with full params (init,cond,post, and body)
		if init == nil || cond == nil || post == nil || body == nil {
			return true
		}

		if *verbose {
			log.Printf("For loop detected (init=%v, cond=%v, post=%v, body=%v)", init, cond, post, body)
		}

		loopInitIdent := ast.NewIdent(fmt.Sprintf("LOOP_INIT_%s", randStringRunes(6)))
		loopCondIdent := ast.NewIdent(fmt.Sprintf("LOOP_COND_%s", randStringRunes(6)))
		loopBodyIdent := ast.NewIdent(fmt.Sprintf("LOOP_BODY_%s", randStringRunes(6)))
		loopEndIdent := ast.NewIdent(fmt.Sprintf("LOOP_END_%s", randStringRunes(6)))

		// Convert `break` in body to `goto loopEndIdent` since
		// after obfuscation it'll be illegal to use `break`
		// outside loop.
		astutil.Apply(body, func(crn *astutil.Cursor) bool {
			branch, ok := crn.Node().(*ast.BranchStmt)
			if !ok {
				return true
			}
			if branch.Tok == token.BREAK {
				crn.Replace(&ast.BranchStmt{
					Tok:   token.GOTO,
					Label: loopEndIdent,
				})
			}
			// continue might find more breaks in the body
			return true
		}, nil)

		body.List = append(body.List, post,
			&ast.BranchStmt{
				Tok:   token.GOTO,
				Label: loopCondIdent,
			})

		loopBegStmts := []ast.Stmt{
			&ast.BranchStmt{
				Tok:   token.GOTO,
				Label: loopInitIdent,
			},
		}

		loopInitStmts := []ast.Stmt{
			&ast.LabeledStmt{
				Label: loopInitIdent,
				Stmt:  &ast.EmptyStmt{},
			},
			init,
			&ast.BranchStmt{
				Tok:   token.GOTO,
				Label: loopCondIdent,
			},
		}

		loopCondStmts := []ast.Stmt{
			&ast.LabeledStmt{
				Label: loopCondIdent,
				Stmt: &ast.IfStmt{
					Cond: cond,
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.BranchStmt{
								Tok:   token.GOTO,
								Label: loopBodyIdent,
							},
						},
					},
					Else: &ast.BranchStmt{
						Tok:   token.GOTO,
						Label: loopEndIdent,
					},
				},
			},
		}

		loopBodyStmts := []ast.Stmt{
			&ast.LabeledStmt{
				Label: loopBodyIdent,
				Stmt:  body,
			},
		}

		loopEndStmts := []ast.Stmt{
			&ast.LabeledStmt{
				Label: loopEndIdent,
				Stmt:  &ast.BlockStmt{},
			},
		}

		ObfuscationBody := []ast.Stmt{}
		ObfuscationBody = append(ObfuscationBody, loopBegStmts...)
		ObfuscationBody = append(ObfuscationBody, loopInitStmts...)
		ObfuscationBody = append(ObfuscationBody, loopCondStmts...)
		ObfuscationBody = append(ObfuscationBody, loopBodyStmts...)
		ObfuscationBody = append(ObfuscationBody, loopEndStmts...)

		cr.Replace(&ast.BlockStmt{
			List: ObfuscationBody,
		})
		return true // replace more fors so we can cover nested ones
	}, nil)
}
