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

// ── Section comment width successes ──────────────────────────────────────────

func TestSectionCommentWidth_SuccessValid(t *testing.T) {
	src := "package p\n\n" + sectionComment("── valid ──") + "\nfunc Foo() {}\n"
	diags := runAnalyzer(t, checkers.SectionCommentWidth, src)
	assert.Empty(t, diags)
}

func TestSectionCommentWidth_SuccessNoSectionComment(t *testing.T) {
	diags := runAnalyzer(t, checkers.SectionCommentWidth, "package p\n\n// a normal comment\nfunc Foo() {}\n")
	assert.Empty(t, diags)
}

// ── Section comment width errors ─────────────────────────────────────────────

func TestSectionCommentWidth_ErrorWrongWidth(t *testing.T) {
	src := "package p\n\n// ── too short ──\nfunc Foo() {}\n"
	diags := runAnalyzer(t, checkers.SectionCommentWidth, src)
	require.Len(t, diags, 1)
	assert.Contains(t, diags[0].Message, "section comment is")
	assert.Contains(t, diags[0].Message, "want 80")
}

// ── Test naming successes ────────────────────────────────────────────────────

func TestTestNaming_SuccessWithSuffix(t *testing.T) {
	src := "package p\n\nimport \"testing\"\n\nfunc TestFoo_Success(t *testing.T) {}\n"
	diags := runAnalyzer(t, checkers.TestNaming, src)
	assert.Empty(t, diags)
}

func TestTestNaming_SuccessWithErrorSuffix(t *testing.T) {
	src := "package p\n\nimport \"testing\"\n\nfunc TestFoo_ErrorNotFound(t *testing.T) {}\n"
	diags := runAnalyzer(t, checkers.TestNaming, src)
	assert.Empty(t, diags)
}

func TestTestNaming_SuccessTestMain(t *testing.T) {
	src := "package p\n\nimport \"testing\"\n\nfunc TestMain(m *testing.M) {}\n"
	diags := runAnalyzer(t, checkers.TestNaming, src)
	assert.Empty(t, diags)
}

func TestTestNaming_SuccessNonTest(t *testing.T) {
	src := "package p\n\nfunc HelperFunc() {}\n"
	diags := runAnalyzer(t, checkers.TestNaming, src)
	assert.Empty(t, diags)
}

func TestTestNaming_SuccessBenchmark(t *testing.T) {
	src := "package p\n\nimport \"testing\"\n\nfunc BenchmarkFoo(b *testing.B) {}\n"
	diags := runAnalyzer(t, checkers.TestNaming, src)
	assert.Empty(t, diags)
}

// ── Test naming errors ───────────────────────────────────────────────────────

func TestTestNaming_ErrorMissingSuffix(t *testing.T) {
	src := "package p\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {}\n"
	diags := runAnalyzer(t, checkers.TestNaming, src)
	require.Len(t, diags, 1)
	assert.Contains(t, diags[0].Message, "must contain _Success or _Error")
}
