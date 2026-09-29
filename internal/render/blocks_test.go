package render

import (
	"strings"
	"testing"

	"github.com/inkcheck/inkpress/internal/deck"
)

func convert(t *testing.T, src string) string {
	t.Helper()
	r := &renderer{deck: &deck.Deck{Dir: t.TempDir()}}
	out, err := r.markdown(&deck.Slide{Path: "slides/x.md"})(src)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func TestDivs(t *testing.T) {
	html := convert(t, "::: box note\n**Bold** text\n\n::: card\n- item\n:::\n:::\n\nAfter")
	for _, want := range []string{
		`<div class="box note">`,
		`<p><strong>Bold</strong> text</p>`,
		`<div class="card">`,
		`<li>item</li>`,
		"</div>\n</div>",
		`<p>After</p>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("HTML lacks %q:\n%s", want, html)
		}
	}
}

func TestDivsLeaveCodeAlone(t *testing.T) {
	html := convert(t, "```\n::: box\n```\n:::")
	if strings.Contains(html, "<div") || !strings.Contains(html, "::: box") {
		t.Errorf("fenced code changed:\n%s", html)
	}
	// A bare ::: with nothing open stays text.
	if html := convert(t, ":::"); !strings.Contains(html, "<p>:::</p>") {
		t.Errorf("stray fence:\n%s", html)
	}
}

func TestDivsCloseAtEnd(t *testing.T) {
	html := convert(t, "::: stat\n# 42%")
	if strings.Count(html, "<div") != 1 || strings.Count(html, "</div>") != 1 {
		t.Errorf("unbalanced:\n%s", html)
	}
}

func TestCallouts(t *testing.T) {
	html := convert(t, "> [!WARNING]\n> Mind the *gap*.")
	for _, want := range []string{
		`<blockquote class="callout warning">`,
		`<p class="callout-title">Warning</p>`,
		`<p>Mind the <em>gap</em>.</p>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("HTML lacks %q:\n%s", want, html)
		}
	}
	if strings.Contains(html, "[!") {
		t.Errorf("marker left in:\n%s", html)
	}

	html = convert(t, "> [!tip] Try this first\n>\n> Body")
	if !strings.Contains(html, `<p class="callout-title">Try this first</p>`) || !strings.Contains(html, "<p>Body</p>") {
		t.Errorf("custom title:\n%s", html)
	}

	html = convert(t, "> Just a quote")
	if strings.Contains(html, "callout") {
		t.Errorf("plain quote marked up:\n%s", html)
	}
}
