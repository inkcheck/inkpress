package deck

import (
	"strings"
	"testing"
)

func TestParseSlideYAML(t *testing.T) {
	s, err := ParseSlide("slides/01.md", "---\nlayout: title\nfooter: false\n---\n\n# Hello\n\nBody\n\n!--\nSay hello.\n")
	if err != nil {
		t.Fatal(err)
	}
	if s.String(KeyLayout) != "title" || s.Bool(KeyFooter, true) {
		t.Errorf("meta = %v", s.Meta)
	}
	if s.Body != "# Hello\n\nBody" {
		t.Errorf("body = %q", s.Body)
	}
	if s.Notes != "Say hello." {
		t.Errorf("notes = %q", s.Notes)
	}
}

func TestParseSlideTOML(t *testing.T) {
	s, err := ParseSlide("x.md", "+++\ntheme = \"dark\"\n+++\nBody")
	if err != nil {
		t.Fatal(err)
	}
	if s.String(KeyTheme) != "dark" || s.Body != "Body" {
		t.Errorf("got %v %q", s.Meta, s.Body)
	}
}

func TestParseSlideNoFrontmatter(t *testing.T) {
	s, err := ParseSlide("x.md", "# Just a slide\r\n")
	if err != nil {
		t.Fatal(err)
	}
	if s.Body != "# Just a slide" || len(s.Meta) != 0 || s.Notes != "" {
		t.Errorf("got %v %q %q", s.Meta, s.Body, s.Notes)
	}
}

func TestParseSlideUnclosed(t *testing.T) {
	if _, err := ParseSlide("x.md", "---\nlayout: x\n# Body"); err == nil {
		t.Error("want an error for unclosed frontmatter")
	}
}

func TestNotesMarkerInCodeIsIgnored(t *testing.T) {
	s, err := ParseSlide("x.md", "```\n!--\n```\n\nText\n!--\nNote")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(s.Body, "```\n!--\n```") || s.Notes != "Note" {
		t.Errorf("body %q notes %q", s.Body, s.Notes)
	}
}

func TestPageSize(t *testing.T) {
	for _, tc := range []struct {
		page Page
		w, h float64
		err  bool
	}{
		{Page{}, 1280, 720, false},
		{Page{Aspect: "4:3", Width: 1024}, 1024, 768, false},
		{Page{Width: 1000, Height: 500}, 1000, 500, false},
		{Page{Aspect: "wide"}, 0, 0, true},
		{Page{Aspect: "16:9", Height: 700}, 0, 0, true},
	} {
		w, h, err := tc.page.Size()
		if (err != nil) != tc.err || w != tc.w || h != tc.h {
			t.Errorf("%+v: got %g×%g, %v", tc.page, w, h, err)
		}
	}
}

func TestSettingsUnknownKey(t *testing.T) {
	if _, err := ParseSettings("titel = \"x\"\n"); err == nil {
		t.Error("want an error for an unknown key")
	}
	s, err := ParseSettings("[meta]\nanything = 1\n")
	if err != nil {
		t.Fatal(err)
	}
	if s.Template != DefaultTemplate {
		t.Errorf("template = %q", s.Template)
	}
}
