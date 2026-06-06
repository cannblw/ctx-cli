package checkers

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

func checkSectionCommentWidth(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		for _, cg := range f.Comments {
			for _, c := range cg.List {
				text := c.Text
				if len(text) > 0 && text[len(text)-1] == '\n' {
					text = text[:len(text)-1]
				}
				if !strings.HasPrefix(text, "// ── ") {
					continue
				}
				width := utf8.RuneCountInString(text)
				if width != 80 {
					pass.Report(analysis.Diagnostic{
						Pos:     c.Pos(),
						Message: fmt.Sprintf("section comment is %d chars, want 80", width),
					})
				}
			}
		}
	}
	return nil, nil
}
