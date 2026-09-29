package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCSS(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "x")
	for name, text := range map[string]string{
		"template.toml":        "",
		"layouts/content.html": "{{.}}",
		"style.css":            "@import url(other.css);\n@font-face { font-family: F; src: url(fonts/f.ttf); }\n.slide { background: url('bg.png') }\n.a { background: url(data:image/png;base64,xx) }",
		"themes/light.css":     ".slide { --fg: black }",
	} {
		p := filepath.Join(tpl, name)
		os.MkdirAll(filepath.Dir(p), 0o755)
		os.WriteFile(p, []byte(text), 0o644)
	}
	tp, err := Load(dir, "x")
	if err != nil {
		t.Fatal(err)
	}
	ss, err := tp.CSS([]string{"light"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"@import", "@font-face", FileURL(filepath.Join(tpl, "fonts", "f.ttf"))} {
		if !strings.Contains(ss.Top, want) {
			t.Errorf("top level lacks %q:\n%s", want, ss.Top)
		}
	}
	for _, want := range []string{"@scope (.tpl-x) {", "@scope (.tpl-x.theme-light) {", FileURL(filepath.Join(tpl, "bg.png")), "url(data:image/png;base64,xx)"} {
		if !strings.Contains(ss.Scoped, want) {
			t.Errorf("scoped lacks %q:\n%s", want, ss.Scoped)
		}
	}
	if strings.Contains(ss.Scoped, "@font-face") {
		t.Error("@font-face left inside @scope")
	}
}

func TestLoadMissing(t *testing.T) {
	if _, err := Load(t.TempDir(), "nope"); err == nil || !strings.Contains(err.Error(), "inkpress new") {
		t.Errorf("err = %v", err)
	}
	if _, err := Load(t.TempDir(), "../up"); err == nil {
		t.Error("want an error for a path as a name")
	}
}
