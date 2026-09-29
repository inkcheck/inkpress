// Package render turns a deck into one HTML document, one page per slide,
// ready for a browser to print.
package render

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/inkcheck/inkpress/internal/deck"
	"github.com/inkcheck/inkpress/internal/templates"
)

// Data is what layouts and partials see.
type Data struct {
	Deck DeckInfo
	// Number counts from 1.
	Number, Total                  int
	Template, Layout, Theme, Title string
	// Heading is the slide's leading # or ## heading, Body the rest, and
	// Content both. Columns is Body split at ||| lines.
	Heading, Body, Content template.HTML
	Columns                []template.HTML
	// Image is the frontmatter image, resolved like a Markdown image.
	Image template.URL
	// Meta is the slide's frontmatter.
	Meta map[string]any
	// Header and Footer are the rendered partials, empty when the slide
	// turns them off.
	Header, Footer template.HTML
}

// DeckInfo is the deck-wide part of Data.
type DeckInfo struct {
	Title, Author, Date string
	Meta                map[string]any
}

// Result is a rendered deck.
type Result struct {
	HTML          []byte
	Width, Height float64
	// Warnings are problems that did not stop the render.
	Warnings []string
}

type renderer struct {
	deck      *deck.Deck
	templates map[string]*templates.Template
	// themes used per template, for the CSS.
	themes   map[string]map[string]bool
	warnings []string
	inkline  inkline
}

// HTML renders the deck.
func HTML(d *deck.Deck) (*Result, error) {
	w, h, err := d.Settings.Page.Size()
	if err != nil {
		return nil, err
	}
	r := &renderer{deck: d, templates: map[string]*templates.Template{}, themes: map[string]map[string]bool{}}
	defer r.inkline.cleanup()

	var slides bytes.Buffer
	for i, s := range d.Slides {
		sec, err := r.slide(s, i+1)
		if err != nil {
			return nil, err
		}
		slides.WriteString(sec)
	}

	css, err := r.css(w, h)
	if err != nil {
		return nil, err
	}
	var doc bytes.Buffer
	fmt.Fprintf(&doc, `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>%s</title>
<base href="%s/">
<style>
%s
</style>
</head>
<body>
%s</body>
</html>
`, html.EscapeString(d.Title()), templates.FileURL(filepath.Join(d.Dir, deck.SlidesDir)), css, slides.String())
	return &Result{HTML: doc.Bytes(), Width: w, Height: h, Warnings: r.warnings}, nil
}

func (r *renderer) warn(format string, args ...any) {
	r.warnings = append(r.warnings, fmt.Sprintf(format, args...))
}

func (r *renderer) template(name string) (*templates.Template, error) {
	if t, ok := r.templates[name]; ok {
		return t, nil
	}
	t, err := templates.Load(filepath.Join(r.deck.Dir, deck.TemplatesDir), name)
	if err != nil {
		return nil, err
	}
	r.templates[name] = t
	r.themes[name] = map[string]bool{}
	return t, nil
}

func (r *renderer) slide(s *deck.Slide, n int) (string, error) {
	set := r.deck.Settings
	name := firstOf(s.String(deck.KeyTemplate), set.Template)
	t, err := r.template(name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", s.Path, err)
	}
	layout := firstOf(s.String(deck.KeyLayout), t.Config.Layout)
	if !t.HasLayout(layout) {
		return "", fmt.Errorf("%s: template %q has no layout %q (layouts: %s)", s.Path, name, layout, strings.Join(t.Layouts(), ", "))
	}
	theme := firstOf(s.String(deck.KeyTheme), set.Theme, t.Config.Theme)
	if !t.HasTheme(theme) {
		return "", fmt.Errorf("%s: template %q has no theme %q (themes: %s)", s.Path, name, theme, strings.Join(t.Themes(), ", "))
	}
	r.themes[name][theme] = true

	body, err := r.diagrams(s, theme)
	if err != nil {
		return "", err
	}
	md := r.markdown(s)
	headingSrc, rest := splitHeading(body)
	lines := strings.Split(rest, "\n")
	var parts []string
	start := 0
	for _, i := range deck.Markers(lines, deck.ColumnMarker) {
		parts = append(parts, strings.Join(lines[start:i], "\n"))
		start = i + 1
	}
	parts = append(parts, strings.Join(lines[start:], "\n"))

	data := Data{
		Deck:     DeckInfo{Title: r.deck.Title(), Author: set.Author, Date: set.Date, Meta: set.Meta},
		Number:   n,
		Total:    len(r.deck.Slides),
		Template: name,
		Layout:   layout,
		Theme:    theme,
		Meta:     s.Meta,
	}
	if data.Heading, err = md(headingSrc); err != nil {
		return "", err
	}
	if data.Body, err = md(strings.Join(parts, "\n\n")); err != nil {
		return "", err
	}
	data.Content = data.Heading + data.Body
	for _, p := range parts {
		c, err := md(p)
		if err != nil {
			return "", err
		}
		data.Columns = append(data.Columns, c)
	}
	data.Title = firstOf(s.String(deck.KeyTitle), headingText(headingSrc))
	if img := s.String(deck.KeyImage); img != "" {
		data.Image = template.URL(r.resolve(s, img))
	}
	if s.Bool(deck.KeyHeader, true) {
		if data.Header, err = t.Partial("header", data); err != nil {
			return "", fmt.Errorf("%s: %w", s.Path, err)
		}
	}
	if s.Bool(deck.KeyFooter, true) {
		if data.Footer, err = t.Partial("footer", data); err != nil {
			return "", fmt.Errorf("%s: %w", s.Path, err)
		}
	}
	inner, err := t.Layout(layout, data)
	if err != nil {
		return "", fmt.Errorf("%s: %w", s.Path, err)
	}

	// The template and theme classes go on the page, since @scope matches
	// only inside its root; the slide is then in scope for .slide rules.
	page := t.Class() + " " + templates.ThemeClass(theme)
	classes := []string{"slide", "layout-" + layout}
	if c := s.String(deck.KeyClass); c != "" {
		classes = append(classes, strings.Fields(c)...)
	}
	style := ""
	if bg := s.String(deck.KeyBackground); bg != "" {
		// A colour, gradient or image; a path is resolved like an image.
		if !strings.ContainsAny(bg, "(#") && looksLikePath(bg) {
			bg = `url("` + r.resolve(s, bg) + `") center / cover no-repeat`
		}
		style = fmt.Sprintf(` style="background: %s"`, html.EscapeString(bg))
	}
	return fmt.Sprintf(`<div class="page %s">
<section class="%s" data-slide="%d" data-source="%s"%s>
%s
</section>
</div>
`, page, html.EscapeString(strings.Join(classes, " ")), n, html.EscapeString(filepath.ToSlash(s.Path)), style, inner), nil
}

var headingRe = regexp.MustCompile(`^#{1,2}\s+(.*?)\s*#*\s*$`)

// splitHeading takes a leading # or ## heading off the body.
func splitHeading(body string) (heading, rest string) {
	first, rest, _ := strings.Cut(body, "\n")
	if headingRe.MatchString(first) {
		return first, strings.TrimLeft(rest, "\n")
	}
	return "", body
}

var inlineMarkRe = regexp.MustCompile("[*_`]")

// headingText is the plain text of a heading line.
func headingText(line string) string {
	m := headingRe.FindStringSubmatch(line)
	if m == nil {
		return ""
	}
	return inlineMarkRe.ReplaceAllString(m[1], "")
}

func looksLikePath(s string) bool {
	ext := strings.ToLower(filepath.Ext(s))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".svg", ".gif", ".webp", ".avif":
		return true
	}
	return false
}

// resolve turns a path from a slide into a URL. Relative paths are tried
// against the slides directory, then the pack, so both ../assets/x.png and
// assets/x.png work.
func (r *renderer) resolve(s *deck.Slide, p string) string {
	if p == "" || strings.Contains(p, "://") || strings.HasPrefix(p, "data:") || strings.HasPrefix(p, "#") {
		return p
	}
	if filepath.IsAbs(p) {
		return templates.FileURL(p)
	}
	fromSlide := filepath.Join(r.deck.Dir, filepath.Dir(s.Path), filepath.FromSlash(p))
	fromPack := filepath.Join(r.deck.Dir, filepath.FromSlash(p))
	for _, c := range []string{fromSlide, fromPack} {
		if _, err := os.Stat(c); err == nil {
			return templates.FileURL(c)
		}
	}
	r.warn("%s: %s not found", s.Path, p)
	return templates.FileURL(fromSlide)
}

func (r *renderer) css(w, h float64) (string, error) {
	var top, scoped strings.Builder
	names := make([]string, 0, len(r.templates))
	for n := range r.templates {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		t := r.templates[n]
		var themes []string
		for th := range r.themes[n] {
			themes = append(themes, th)
		}
		sort.Strings(themes)
		ss, err := t.CSS(themes)
		if err != nil {
			return "", err
		}
		top.WriteString(ss.Top)
		scoped.WriteString(ss.Scoped)
		for _, th := range themes {
			code, err := codeCSS(t.CodeStyle(th))
			if err != nil {
				return "", fmt.Errorf("template %q, theme %q: %w", n, th, err)
			}
			fmt.Fprintf(&scoped, "@scope (.%s.%s) {\n%s}\n", t.Class(), templates.ThemeClass(th), code)
		}
	}
	return top.String() + fmt.Sprintf(baseCSS, w, h, w, h) + scoped.String(), nil
}

// baseCSS sizes the pages. Each slide sits in a .page that is a size
// container, so templates can size things in cqw and cqh units.
const baseCSS = `@page { size: %[1]gpx %[2]gpx; margin: 0; }
:root { --slide-width: %[3]gpx; --slide-height: %[4]gpx; }
*, *::before, *::after { box-sizing: border-box; }
html, body { margin: 0; padding: 0; }
.page {
  width: var(--slide-width); height: var(--slide-height);
  position: relative; overflow: hidden;
  container-type: size;
  break-after: page; break-inside: avoid;
}
.page:last-child { break-after: auto; }
.slide {
  width: 100%%; height: 100%%; position: relative; overflow: hidden;
  -webkit-print-color-adjust: exact; print-color-adjust: exact;
}
.slide img { max-width: 100%%; }
@media screen {
  body { background: #2b2b2b; padding: 24px 0; }
  .page { margin: 0 auto 24px; box-shadow: 0 4px 24px rgb(0 0 0 / 0.5); }
}
`

func firstOf(s ...string) string {
	for _, x := range s {
		if x != "" {
			return x
		}
	}
	return ""
}
