package checkers

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func checkTestNaming(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		if !strings.HasSuffix(pass.Fset.File(f.Pos()).Name(), "_test.go") {
			continue
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			name := fn.Name.Name
			if !strings.HasPrefix(name, "Test") {
				continue
			}
			if strings.HasPrefix(name, "TestMain") {
				continue
			}
			if strings.Contains(name, "_Success") || strings.Contains(name, "_Error") {
				continue
			}
			pass.Report(analysis.Diagnostic{
				Pos:     fn.Pos(),
				Message: fmt.Sprintf("test function %q must contain _Success or _Error in its name", name),
			})
		}
	}
	return nil, nil
}
