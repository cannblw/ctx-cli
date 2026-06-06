package checkers

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func Main() {
	unitchecker.Main(
		SectionCommentWidth,
		TestNaming,
	)
}

var SectionCommentWidth = &analysis.Analyzer{
	Name: "sectioncommentwidth",
	Doc:  "enforces section comments (// ── ...) are exactly 80 characters wide",
	Run:  checkSectionCommentWidth,
}

var TestNaming = &analysis.Analyzer{
	Name: "testnaming",
	Doc:  "enforces every test function name contains _Success or _Error",
	Run:  checkTestNaming,
}
