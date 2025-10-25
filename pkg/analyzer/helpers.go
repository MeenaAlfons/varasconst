package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
)

func visitGlobalDeclarations(ins *inspector.Inspector, visit func(*ast.GenDecl)) {
	ins.Nodes([]ast.Node{
		// Global variables are of type GenDecl.
		(*ast.GenDecl)(nil),

		// Visit function declarations to skip their bodies
		// in order to not visit local variable declarations.
		(*ast.FuncDecl)(nil),
	}, func(node ast.Node, push bool) (proceed bool) {
		_, isFunc := node.(*ast.FuncDecl)
		if isFunc {
			return false
		}

		genDecl := node.(*ast.GenDecl)
		visit(genDecl)
		return true
	})
}

func filterVars(visit func(*ast.GenDecl)) func(*ast.GenDecl) {
	return func(genDecl *ast.GenDecl) {
		if genDecl.Tok != token.VAR {
			return
		}
		visit(genDecl)
	}
}

func filterMarkedConst(visit func(*ast.Ident)) func(*ast.GenDecl) {
	return func(genDecl *ast.GenDecl) {
		if genDecl.Doc != nil &&
			strings.TrimSpace(genDecl.Doc.Text()) == "const" {
			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range valueSpec.Names {
					visit(name)
				}
			}
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			if valueSpec.Doc != nil && strings.TrimSpace(valueSpec.Doc.Text()) == "const" {
				for _, name := range valueSpec.Names {
					visit(name)
				}
			}
		}
	}
}

func visitVarsMarkedConst(ins *inspector.Inspector, visit func(*ast.Ident)) {
	visitGlobalDeclarations(ins,
		filterVars(
			filterMarkedConst(
				visit,
			),
		))
}
