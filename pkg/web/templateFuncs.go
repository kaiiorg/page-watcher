package web

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"html"
	"html/template"
	"strings"
)

func (w *WebDiffPreviewer) replaceEndLines(s string) template.HTML {
	s = strings.ReplaceAll(
		html.EscapeString(s),
		"\n",
		"<br>",
	)
	return template.HTML(s)
}

func (w *WebDiffPreviewer) isInsert(diff diffmatchpatch.Diff) bool {
	return diff.Type == diffmatchpatch.DiffInsert
}

func (w *WebDiffPreviewer) isDel(diff diffmatchpatch.Diff) bool {
	return diff.Type == diffmatchpatch.DiffDelete
}

func (w *WebDiffPreviewer) isEqual(diff diffmatchpatch.Diff) bool {
	return diff.Type == diffmatchpatch.DiffEqual
}
