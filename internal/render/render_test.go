package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/inkcheck/inkpress/internal/deck"
	"github.com/inkcheck/inkpress/internal/scaffold"
)

func samplePack(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if _, err := scaffold.New(dir, false); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestSamplePack(t *testing.T) {
	// No inkline: diagram blocks stay code, with a warning.
	t.Setenv("PATH", "")
	t.Setenv(InklineEnv, "")
	d, err := deck.Load(samplePack(t))
	if err != nil {
		t.Fatal(err)
	}
	res, err := HTML(d)
	if err != nil {
		t.Fatal(err)
	}
	html := string(res.HTML)
	if res.Width != 1280 || res.Height != 720 {
		t.Errorf("size %g×%g", res.Width, res.Height)
	}
	for _, want := range []string{
		`<div class="page tpl-default theme-light">`,
		`<div class="page tpl-default theme-dark">`,
		`class="slide layout-two-column"`,
		`<div class="column">`,
		`@scope (.tpl-default.theme-dark)`,
		`assets/chart.svg`,
		`class="chroma"`,
		`2 / 8`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("HTML lacks %q", want)
		}
	}
	if strings.Contains(html, "Notes go after") {
		t.Error("speaker notes leaked into the HTML")
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "inkline not found") {
		t.Errorf("warnings = %q", res.Warnings)
	}
}

func TestErrors(t *testing.T) {
	dir := samplePack(t)
	write := func(name, text string) {
		os.WriteFile(filepath.Join(dir, "slides", name), []byte(text), 0o644)
	}
	write("99-bad.md", "---\nlayout: nope\n---\n# x")
	d, err := deck.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := HTML(d); err == nil || !strings.Contains(err.Error(), `no layout "nope"`) {
		t.Errorf("err = %v", err)
	}
	write("99-bad.md", "---\ntheme: sepia\n---\n# x")
	d, _ = deck.Load(dir)
	if _, err := HTML(d); err == nil || !strings.Contains(err.Error(), `no theme "sepia"`) {
		t.Errorf("err = %v", err)
	}
	write("99-bad.md", "---\nskip: true\ntheme: sepia\n---\n# x")
	d, _ = deck.Load(dir)
	if _, err := HTML(d); err != nil {
		t.Errorf("skipped slide still rendered: %v", err)
	}
}

func TestSplitHeading(t *testing.T) {
	h, rest := splitHeading("## Title *here*\n\nBody")
	if h != "## Title *here*" || rest != "Body" || headingText(h) != "Title here" {
		t.Errorf("%q %q %q", h, rest, headingText(h))
	}
	if h, _ := splitHeading("Text\n# Late"); h != "" {
		t.Errorf("heading %q", h)
	}
}
