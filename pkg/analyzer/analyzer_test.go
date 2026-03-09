package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AntiSlang/linter/pkg/analyzer"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	wd, _ := os.Getwd()
	testdata := filepath.Join(wd, "testdata")

	analysistest.Run(t, testdata, analyzer.Analyzer, "linter_test")
}
