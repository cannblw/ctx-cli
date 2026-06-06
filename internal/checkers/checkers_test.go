package checkers_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"

	"github.com/cannblw/ctx-cli/internal/checkers"
)

func runAnalyzer(t *testing.T, a *analysis.Analyzer, src string) []analysis.Diagnostic {
	t.Helper()
	var diags []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: a,
		Fset:     token.NewFileSet(),
		Report: func(d analysis.Diagnostic) {
			diags = append(diags, d)
		},
	}

	f, err := parser.ParseFile(pass.Fset, "test_test.go", src, parser.ParseComments)
	require.NoError(t, err)
	pass.Files = []*ast.File{f}

	_, err = a.Run(pass)
	require.NoError(t, err)
	return diags
}

func sectionComment(text string) string {
	prefix := "// "
	visual := utf8.RuneCountInString(prefix) + utf8.RuneCountInString(text)
	pad := strings.Repeat("─", 80-visual)
	return prefix + text + pad
}

func TestSectionCommentWidth_Success(t *testing.T) {
	t.Run("valid 80-char comment passes", func(t *testing.T) {
		src := "package p\n\n" + sectionComment("── valid ──") + "\nfunc Foo() {}\n"
		diags := runAnalyzer(t, checkers.SectionCommentWidth, src)
		assert.Empty(t, diags)
	})

	t.Run("no section comment passes", func(t *testing.T) {
		diags := runAnalyzer(t, checkers.SectionCommentWidth, "package p\n\n// a normal comment\nfunc Foo() {}\n")
		assert.Empty(t, diags)
	})

	t.Run("wrong width reports diagnostic", func(t *testing.T) {
		src := "package p\n\n// ── too short ──\nfunc Foo() {}\n"
		diags := runAnalyzer(t, checkers.SectionCommentWidth, src)
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "section comment is")
	})
}

func TestTestNaming_Success(t *testing.T) {
	t.Run("test with _Success passes", func(t *testing.T) {
		src := "package p\n\nimport \"testing\"\n\nfunc TestFoo_Success(t *testing.T) {}\n"
		diags := runAnalyzer(t, checkers.TestNaming, src)
		assert.Empty(t, diags)
	})

	t.Run("test with _Error passes", func(t *testing.T) {
		src := "package p\n\nimport \"testing\"\n\nfunc TestFoo_ErrorNotFound(t *testing.T) {}\n"
		diags := runAnalyzer(t, checkers.TestNaming, src)
		assert.Empty(t, diags)
	})

	t.Run("TestMain passes without suffix", func(t *testing.T) {
		src := "package p\n\nimport \"testing\"\n\nfunc TestMain(m *testing.M) {}\n"
		diags := runAnalyzer(t, checkers.TestNaming, src)
		assert.Empty(t, diags)
	})

	t.Run("non-test function passes", func(t *testing.T) {
		src := "package p\n\nfunc HelperFunc() {}\n"
		diags := runAnalyzer(t, checkers.TestNaming, src)
		assert.Empty(t, diags)
	})

	t.Run("benchmark passes without suffix", func(t *testing.T) {
		src := "package p\n\nimport \"testing\"\n\nfunc BenchmarkFoo(b *testing.B) {}\n"
		diags := runAnalyzer(t, checkers.TestNaming, src)
		assert.Empty(t, diags)
	})

	t.Run("test without suffix reports diagnostic", func(t *testing.T) {
		src := "package p\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {}\n"
		diags := runAnalyzer(t, checkers.TestNaming, src)
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "must contain _Success or _Error")
	})
}
