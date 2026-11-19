package inkdeck

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/inkcheck/inkpress/scaffold"
)

func TestShare(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "talk")
	if _, err := scaffold.New(dir, false); err != nil {
		t.Fatal(err)
	}
	// Things that stay out of a shared deck.
	for name, body := range map[string]string{
		"talk.pdf": "pdf", "talk.pptx": "x", "talk.html": "x", "old.inkdeck": "x",
		".git/HEAD": "ref", "slides/.inkpress-x.d2": "x", ".DS_Store": "x",
	} {
		os.MkdirAll(filepath.Join(dir, filepath.Dir(name)), 0o755)
		os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644)
	}
	os.Symlink("/etc/hosts", filepath.Join(dir, "assets", "link"))

	// At the top of the pack is fine: it leaves itself out.
	out := filepath.Join(dir, "talk.inkdeck")
	n, err := Share(dir, out)
	if err != nil {
		t.Fatal(err)
	}
	zr, err := zip.OpenReader(out)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	names := map[string]bool{}
	for _, f := range zr.File {
		names[f.Name] = true
	}
	if len(names) != n {
		t.Errorf("%d entries, %d reported", len(names), n)
	}
	for _, want := range []string{"talk/settings.toml", "talk/slides/01-title.md", "talk/.templates/default/template.toml"} {
		if !names[want] {
			t.Errorf("missing %s", want)
		}
	}
	for name := range names {
		for _, bad := range []string{"talk/talk.", "talk/old.inkdeck", ".git/", ".inkpress-", ".DS_Store", "assets/link"} {
			if strings.Contains(name, bad) {
				t.Errorf("packed %s", name)
			}
		}
	}
}

func TestShareRefuses(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "talk")
	if _, err := Share(dir, filepath.Join(t.TempDir(), "x.inkdeck")); err == nil {
		t.Error("shared a folder without slides/")
	}
	if _, err := scaffold.New(dir, false); err != nil {
		t.Fatal(err)
	}
	if _, err := Share(dir, filepath.Join(dir, "slides", "x.inkdeck")); err == nil {
		t.Error("wrote into slides/")
	}
	if _, err := os.Stat(filepath.Join(dir, "slides", "x.inkdeck")); err == nil {
		t.Error("left a file in slides/")
	}
}
