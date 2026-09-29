package render

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/inkcheck/inkpress/internal/deck"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Two kinds of block that templates style: fenced divs and callouts.
//
// A fenced div is Markdown between ::: lines. "::: box note" opens
// <div class="box note"> and a bare ":::" closes it. Divs nest, and a
// div left open closes at the end of the slide or column.
//
// A callout is a blockquote that starts with a [!KIND] line, as on GitHub:
// "> [!NOTE]" becomes <blockquote class="callout note"> with a title
// paragraph. Text after the marker replaces the title.

var (
	divOpenRe  = regexp.MustCompile(`^:{3,}\s*([A-Za-z][\w-]*(?:\s+[A-Za-z][\w-]*)*)\s*:*$`)
	divCloseRe = regexp.MustCompile(`^:{3,}$`)
)

// divs turns ::: fences into HTML divs, leaving fenced code alone. The
// blank lines around each tag end the HTML block, so the Markdown inside
// is still parsed.
func divs(src string) string {
	if !strings.Contains(src, ":::") {
		return src
	}
	var out []string
	fence := ""
	open := 0
	for _, l := range strings.Split(src, "\n") {
		t := strings.TrimSpace(l)
		switch {
		case fence != "":
			if strings.HasPrefix(t, fence) && strings.Trim(t, fence[:1]) == "" {
				fence = ""
			}
		case deck.FenceOpen(t) != "":
			fence = deck.FenceOpen(t)
		case divCloseRe.MatchString(t) && open > 0:
			open--
			out = append(out, "", "</div>", "")
			continue
		case divOpenRe.MatchString(t):
			open++
			out = append(out, "", `<div class="`+strings.Join(strings.Fields(divOpenRe.FindStringSubmatch(t)[1]), " ")+`">`, "")
			continue
		}
		out = append(out, l)
	}
	for ; open > 0; open-- {
		out = append(out, "", "</div>")
	}
	return strings.Join(out, "\n")
}

var calloutRe = regexp.MustCompile(`^\[!([A-Za-z]+)\]\s*(.*)$`)

// callouts marks up blockquotes that start with a [!KIND] line.
type callouts struct{}

func (callouts) Transform(doc *ast.Document, reader text.Reader, _ parser.Context) {
	src := reader.Source()
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		bq, ok := n.(*ast.Blockquote)
		if !ok || !entering {
			return ast.WalkContinue, nil
		}
		p, ok := bq.FirstChild().(*ast.Paragraph)
		if !ok || p.Lines().Len() == 0 {
			return ast.WalkContinue, nil
		}
		first := p.Lines().At(0)
		m := calloutRe.FindSubmatch(bytes.TrimSpace(first.Value(src)))
		if m == nil {
			return ast.WalkContinue, nil
		}
		// Drop the marker line's inline nodes.
		for c := p.FirstChild(); c != nil; {
			next := c.NextSibling()
			t, isText := c.(*ast.Text)
			p.RemoveChild(p, c)
			if isText && (t.SoftLineBreak() || t.HardLineBreak()) {
				break
			}
			c = next
		}
		if p.ChildCount() == 0 {
			bq.RemoveChild(bq, p)
		}
		kind := strings.ToLower(string(m[1]))
		title := string(m[2])
		if title == "" {
			title = strings.ToUpper(kind[:1]) + kind[1:]
		}
		bq.SetAttributeString("class", []byte("callout "+kind))
		tp := ast.NewParagraph()
		tp.SetAttributeString("class", []byte("callout-title"))
		tp.AppendChild(tp, ast.NewString([]byte(title)))
		bq.InsertBefore(bq, bq.FirstChild(), tp)
		return ast.WalkContinue, nil
	})
}
