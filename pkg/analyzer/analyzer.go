package analyzer

import (
	"go/ast"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	reOnlyEnglishAndSpaces = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	sensitiveKeywords      = []string{"password", "token", "api key"}
	logMethods             = map[string]bool{"Debug": true, "Info": true, "Warn": true, "Error": true}
)

var Analyzer = &analysis.Analyzer{
	Name:     "linter",
	Doc:      "checks logging messages for style and security standards",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspectRes := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	inspectRes.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		if !isLoggingCall(pass, call) {
			return
		}

		for _, arg := range call.Args {
			checkSensitiveData(pass, arg)
		}

		if len(call.Args) > 0 {
			if lit, ok := call.Args[0].(*ast.BasicLit); ok {
				msg := strings.Trim(lit.Value, `"`+"`")
				if msg == "" {
					return
				}

				firstRune, _ := utf8.DecodeRuneInString(msg)
				if unicode.IsUpper(firstRune) {
					pass.Reportf(lit.Pos(), "log message should start with a lowercase letter")
				}

				if containsCyrillic(msg) {
					pass.Reportf(lit.Pos(), "log message should be in English only")
					return
				}

				if !reOnlyEnglishAndSpaces.MatchString(msg) {
					pass.Reportf(lit.Pos(), "log message contains forbidden symbols or emojis")
				}
			}
		}
	})

	return nil, nil
}

func isLoggingCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if !logMethods[selector.Sel.Name] {
		return false
	}
	if selection, ok := pass.TypesInfo.Selections[selector]; ok {
		pkg := selection.Obj().Pkg()
		if pkg != nil {
			path := pkg.Path()
			return path == "log/slog" || path == "go.uber.org/zap" || path == "linter_test"
		}
	}
	return false
}

func checkSensitiveData(pass *analysis.Pass, expr ast.Expr) {
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BasicLit:
			val := strings.ToLower(x.Value)
			for _, key := range sensitiveKeywords {
				if strings.Contains(val, key) {
					pass.Reportf(x.Pos(), "potential sensitive data exposure: %s", key)
				}
			}
		case *ast.Ident:
			name := strings.ToLower(x.Name)
			for _, key := range sensitiveKeywords {
				if strings.Contains(name, key) {
					pass.Reportf(x.Pos(), "avoid logging sensitive variables like '%s'", x.Name)
				}
			}
		}
		return true
	})
}

func containsCyrillic(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Cyrillic, r) {
			return true
		}
	}
	return false
}
