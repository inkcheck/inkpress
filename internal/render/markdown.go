package render

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/inkcheck/inkpress/internal/deck"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// markdown returns a converter for one slide's Markdown.
func (r *renderer) markdown(s *deck.Slide) func(string) (template.HTML, error) {
	gm := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(html.WithClasses(true)),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.Prioritized(&images{r: r, s: s}, 100),
				util.Prioritized(callouts{}, 200),
			),
		),
		// Slides are the author's own; raw HTML is allowed.
		goldmark.WithRendererOptions(gmhtml.WithUnsafe()),
	)
	return func(src string) (template.HTML, error) {
		if strings.TrimSpace(src) == "" {
			return "", nil
		}
		var b bytes.Buffer
		if err := gm.Convert([]byte(divs(src)), &b); err != nil {
			return "", fmt.Errorf("%s: %w", s.Path, err)
		}
		return template.HTML(b.String()), nil
	}
}

// images resolves image paths against the slide and the pack.
type images struct {
	r *renderer
	s *deck.Slide
}

func (t *images) Transform(doc *ast.Document, _ text.Reader, _ parser.Context) {
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if img, ok := n.(*ast.Image); ok && entering {
			img.Destination = []byte(t.r.resolve(t.s, string(img.Destination)))
		}
		return ast.WalkContinue, nil
	})
}

// codeCSS is the stylesheet for Chroma style name.
func codeCSS(name string) (string, error) {
	style := styles.Get(name)
	if style == nil || (style == styles.Fallback && name != "swapoff") {
		return "", fmt.Errorf("unknown code style %q", name)
	}
	var b bytes.Buffer
	if err := html.New(html.WithClasses(true)).WriteCSS(&b, style); err != nil {
		return "", err
	}
	return b.String(), nil
}
