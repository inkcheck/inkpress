// Package templates loads slide template packs from a pack's .templates/.
//
// A template pack is a directory:
//
//	template.toml        description, default layout and theme, code styles
//	style.css            styles for .slide and everything in it
//	themes/<theme>.css   custom properties for .slide, one file per theme
//	layouts/<name>.html  one html/template per layout
//	partials/<name>.html header.html, footer.html and any others
//	assets/, fonts/      files the CSS and layouts refer to
package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config is a template pack's template.toml.
type Config struct {
	Description string `toml:"description"`
	// Layout is used by slides that name none.
	Layout string `toml:"layout"`
	// Theme is used when neither the slide nor the settings name one.
	Theme string `toml:"theme"`
	// CodeStyles maps a theme to a Chroma style for code blocks.
	// See https://xyproto.github.io/splash/docs/.
	CodeStyles map[string]string `toml:"code_styles"`
}

// Template is a loaded template pack.
type Template struct {
	Name   string
	Dir    string
	Config Config
	set    *template.Template
}

// Template names inside the set.
const (
	layoutPrefix  = "layout/"
	partialPrefix = "partial/"
)

// Load reads the template pack name from dir/.templates.
func Load(templatesDir, name string) (*Template, error) {
	if name == "" || strings.ContainsAny(name, `/\`) || strings.HasPrefix(name, ".") {
		return nil, fmt.Errorf("template %q: not a template name", name)
	}
	dir := filepath.Join(templatesDir, name)
	t := &Template{Name: name, Dir: dir, Config: Config{Layout: "content", Theme: "light"}}
	text, err := os.ReadFile(filepath.Join(dir, "template.toml"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("template %q: %s not found (run `inkpress new` to get the default template)", name, filepath.Join(dir, "template.toml"))
		}
		return nil, err
	}
	if _, err := toml.Decode(string(text), &t.Config); err != nil {
		return nil, fmt.Errorf("template %q: template.toml: %w", name, err)
	}

	t.set = template.New(name).Funcs(template.FuncMap{
		"asset": t.Asset,
	})
	for _, kind := range []struct{ dir, prefix string }{{"partials", partialPrefix}, {"layouts", layoutPrefix}} {
		files, _ := filepath.Glob(filepath.Join(dir, kind.dir, "*.html"))
		for _, f := range files {
			src, err := os.ReadFile(f)
			if err != nil {
				return nil, err
			}
			n := kind.prefix + strings.TrimSuffix(filepath.Base(f), ".html")
			if _, err := t.set.New(n).Parse(string(src)); err != nil {
				return nil, fmt.Errorf("template %q: %w", name, err)
			}
		}
	}
	if !t.HasLayout(t.Config.Layout) {
		return nil, fmt.Errorf("template %q: default layout %q has no layouts/%s.html", name, t.Config.Layout, t.Config.Layout)
	}
	return t, nil
}

// Layouts lists the layout names.
func (t *Template) Layouts() []string { return t.names(layoutPrefix) }

// HasLayout reports whether the template has layout name.
func (t *Template) HasLayout(name string) bool { return t.set.Lookup(layoutPrefix+name) != nil }

// Themes lists the theme names.
func (t *Template) Themes() []string {
	files, _ := filepath.Glob(filepath.Join(t.Dir, "themes", "*.css"))
	var names []string
	for _, f := range files {
		names = append(names, strings.TrimSuffix(filepath.Base(f), ".css"))
	}
	return names
}

// HasTheme reports whether the template has theme name.
func (t *Template) HasTheme(name string) bool {
	_, err := os.Stat(filepath.Join(t.Dir, "themes", name+".css"))
	return err == nil
}

func (t *Template) names(prefix string) []string {
	var names []string
	for _, x := range t.set.Templates() {
		if n, ok := strings.CutPrefix(x.Name(), prefix); ok {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	return names
}

// Layout executes layout name with data.
func (t *Template) Layout(name string, data any) (template.HTML, error) {
	return t.exec(layoutPrefix+name, data)
}

// Partial executes partial name with data. A missing partial is empty.
func (t *Template) Partial(name string, data any) (template.HTML, error) {
	if t.set.Lookup(partialPrefix+name) == nil {
		return "", nil
	}
	return t.exec(partialPrefix+name, data)
}

func (t *Template) exec(name string, data any) (template.HTML, error) {
	var b bytes.Buffer
	if err := t.set.ExecuteTemplate(&b, name, data); err != nil {
		return "", fmt.Errorf("template %q: %w", t.Name, err)
	}
	return template.HTML(b.String()), nil
}

// Asset returns the URL of a file in the template's assets/ directory.
func (t *Template) Asset(name string) template.URL {
	return template.URL(FileURL(filepath.Join(t.Dir, "assets", filepath.FromSlash(name))))
}

// CodeStyle returns the Chroma style for theme.
func (t *Template) CodeStyle(theme string) string {
	if s := t.Config.CodeStyles[theme]; s != "" {
		return s
	}
	return "github"
}

// FileURL returns the file: URL of an absolute path.
func FileURL(path string) string {
	p := filepath.ToSlash(path)
	if runtime.GOOS == "windows" {
		p = "/" + p
	}
	return (&url.URL{Scheme: "file", Path: p}).String()
}
