package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:      "varasconst",
	Doc:       "Checks that global variables marked with `// const` are never changed.",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&VarMarkedConstFact{}},
}

func run(pass *analysis.Pass) (interface{}, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Export facts for global variables marked with `// const`
	visitVarsMarkedConst(ins, func(ident *ast.Ident) {
		if pass.TypesInfo == nil {
			// just a protection against TypesInfo nil
			return
		}
		varObj := pass.TypesInfo.ObjectOf(ident)
		if varObj == nil {
			return
		}

		pass.ExportObjectFact(origin(varObj), &VarMarkedConstFact{})
	})

	visitAssignment := func(node *ast.AssignStmt) {
		// For assignment (=) Check for violations for each lhs
		if node.Tok == token.ASSIGN {
			for _, lhs := range node.Lhs {
				checkForViolations(pass, lhs)
			}
		}
	}
	onlyAssignment := func(visit func(*ast.AssignStmt)) func(node ast.Node) {
		return func(node ast.Node) {
			switch node := node.(type) {
			case *ast.AssignStmt:
				visit(node)
			}
		}
	}
	ins.Preorder([]ast.Node{(*ast.AssignStmt)(nil)}, onlyAssignment(visitAssignment))
	return nil, nil
}

// Checks for violations in lhs = ...
func checkForViolations(pass *analysis.Pass, lhs ast.Expr) {
	selector, isSelector := lhs.(*ast.SelectorExpr)
	if isSelector {
		checkForViolationsForSelector(pass, selector)
	}

	// If it's not a selector, it must be a simple identifier
	ident, ok := lhs.(*ast.Ident)
	if !ok {
		return
	}
	if pass.TypesInfo == nil {
		// just a protection against TypesInfo nil
		return
	}
	obj := pass.TypesInfo.ObjectOf(ident)
	if obj != nil && pass.ImportObjectFact(origin(obj), &VarMarkedConstFact{}) {
		pass.Reportf(ident.Pos(), "assignment to global variable marked with const: %s", ident.Name)
	}
}

// Checks for violations in X.Sel = ...
func checkForViolationsForSelector(pass *analysis.Pass, selector *ast.SelectorExpr) {
	// Selector has the form X.Sel
	// X could be another selector.
	// We currently only check for one level selection assuming X is an identifier
	firstPartIdent, ok := selector.X.(*ast.Ident)
	if !ok {
		return
	}

	if pass.TypesInfo == nil {
		// just a protection against TypesInfo nil
		return
	}
	// Check if X is a package name
	if _, isPkgName := pass.TypesInfo.ObjectOf(firstPartIdent).(*types.PkgName); !isPkgName {
		// Assignment to a field in local variable or global variable
		// ...
		return
	}

	// X is definitely a package, let's check if the identifier X.Sel was registered as a VarMarkedConstFact
	if pass.ImportObjectFact(origin(pass.TypesInfo.ObjectOf(selector.Sel)), &VarMarkedConstFact{}) {
		pass.Reportf(selector.Pos(), "assignment to global variable marked with const in another package: %s.%s",
			firstPartIdent.Name, selector.Sel.Name)
	}
}

func origin(obj types.Object) types.Object {
	switch obj := obj.(type) {
	case *types.Func:
		return obj.Origin()
	case *types.Var:
		return obj.Origin()
	}
	return obj
}
