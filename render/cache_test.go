package render

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/inkcheck/inkpress/deck"
)

// fakeInkline writes a script that renders every diagram as the same SVG
// and counts its runs in a file.
func fakeInkline(t *testing.T) (count func() int) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("needs a shell script")
	}
	dir := t.TempDir()
	log := filepath.Join(dir, "runs")
	script := filepath.Join(dir, "inkline")
	src := "#!/bin/sh\necho run >> " + log + "\nwhile [ $# -gt 0 ]; do if [ \"$1\" = -o ]; then out=$2; fi; shift; done\necho '<svg xmlns=\"http://www.w3.org/2000/svg\"/>' > \"$out\"\n"
	if err := os.WriteFile(script, []byte(src), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(InklineEnv, script)
	return func() int {
		b, _ := os.ReadFile(log)
		return strings.Count(string(b), "run")
	}
}

func TestDiagramCache(t *testing.T) {
	runs := fakeInkline(t)
	d, err := deck.Load(samplePack(t))
	if err != nil {
		t.Fatal(err)
	}
	cache := NewDiagramCache()
	if _, err := HTMLWith(d, Options{Diagrams: cache}); err != nil {
		t.Fatal(err)
	}
	first := runs()
	if first == 0 || cache.Len() != first {
		t.Fatalf("runs = %d, cached = %d", first, cache.Len())
	}
	if _, err := HTMLWith(d, Options{Diagrams: cache}); err != nil {
		t.Fatal(err)
	}
	if runs() != first {
		t.Errorf("inkline ran again for unchanged diagrams: %d runs", runs())
	}

	// A changed diagram is rendered again, and the old one is dropped.
	for _, s := range d.Slides {
		s.Body = strings.Replace(s.Body, "```inkline", "```inkline\n# changed", 1)
	}
	if _, err := HTMLWith(d, Options{Diagrams: cache}); err != nil {
		t.Fatal(err)
	}
	if runs() != 2*first {
		t.Errorf("runs = %d, want %d", runs(), 2*first)
	}
	if cache.Len() != first {
		t.Errorf("cached = %d, want %d", cache.Len(), first)
	}
}

func TestSlideResults(t *testing.T) {
	t.Setenv("PATH", "")
	d, err := deck.Load(samplePack(t))
	if err != nil {
		t.Fatal(err)
	}
	res, err := HTML(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Slides) != len(d.Slides) {
		t.Fatalf("%d results for %d slides", len(res.Slides), len(d.Slides))
	}
	first := res.Slides[0]
	if first.Number != 1 || first.Source != d.Slides[0].Path || first.Layout != "title" || first.Template != "default" || first.Theme == "" {
		t.Errorf("slide 1 = %+v", first)
	}
}
