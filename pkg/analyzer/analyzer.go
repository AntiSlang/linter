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
	logMethods             = map[string]bool{"Debug": true, "Info": true, "Warn": true, "Error": true}

	checkLowercase bool
	checkEnglish   bool
	checkSpecials  bool
	checkSensitive bool
	sensitiveWords string
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "checks logging messages for style and security standards",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func init() {
	Analyzer.Flags.BoolVar(&checkLowercase, "check-lowercase", true, "check if log message starts with a lowercase letter")
	Analyzer.Flags.BoolVar(&checkEnglish, "check-english", true, "check if log message is in English only")
	Analyzer.Flags.BoolVar(&checkSpecials, "check-specials", true, "check if log message contains forbidden symbols")
	Analyzer.Flags.BoolVar(&checkSensitive, "check-sensitive", true, "check for sensitive data in logs")
	Analyzer.Flags.StringVar(&sensitiveWords, "sensitive-words", "password,token,api key", "comma-separated list of sensitive words")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspectRes := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	inspectRes.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		if !isLoggingCall(pass, call) {
			return
		}

		keywords := strings.Split(sensitiveWords, ",")
		for i := range keywords {
			keywords[i] = strings.TrimSpace(keywords[i])
		}

		if checkSensitive {
			for _, arg := range call.Args {
				checkSensitiveData(pass, arg, keywords)
			}
		}

		if len(call.Args) > 0 {
			if lit, ok := call.Args[0].(*ast.BasicLit); ok {
				msg := strings.Trim(lit.Value, `"`+"`")
				if msg == "" {
					return
				}

				if checkLowercase {
					firstRune, _ := utf8.DecodeRuneInString(msg)
					if unicode.IsUpper(firstRune) {
						pass.Reportf(lit.Pos(), "log message should start with a lowercase letter")
					}
				}

				if checkEnglish && containsCyrillic(msg) {
					pass.Reportf(lit.Pos(), "log message should be in English only")
					return
				}

				if checkSpecials && !reOnlyEnglishAndSpaces.MatchString(msg) {
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

func checkSensitiveData(pass *analysis.Pass, expr ast.Expr, sensitiveKeywords []string) {
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
