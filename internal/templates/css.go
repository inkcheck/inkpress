package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// A deck can mix templates and themes, so inkpress scopes each stylesheet to
// the slides that use it:
//
//	style.css          @scope (.tpl-<name>) { … }
//	themes/<theme>.css @scope (.tpl-<name>.theme-<theme>) { … }
//
// @font-face and @import cannot sit inside @scope, so they are hoisted to the
// top level. Relative url()s are made absolute, since the CSS is inlined.

var (
	urlRe      = regexp.MustCompile(`url\(\s*(['"]?)([^'")]+)(['"]?)\s*\)`)
	fontFaceRe = regexp.MustCompile(`(?s)@font-face\s*\{[^}]*\}`)
	importRe   = regexp.MustCompile(`(?m)^\s*@(?:import|charset)[^;]*;`)
)

// Class is the CSS class for slides that use the template.
func (t *Template) Class() string { return "tpl-" + t.Name }

// ThemeClass is the CSS class for slides that use theme.
func ThemeClass(theme string) string { return "theme-" + theme }

// Stylesheet is CSS split into rules that must be at the top level and rules
// that can be scoped.
type Stylesheet struct {
	Top, Scoped string
}

// CSS returns the template's style.css and the named themes, scoped.
func (t *Template) CSS(themes []string) (Stylesheet, error) {
	var top, scoped strings.Builder
	add := func(file, scope string) error {
		src, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("template %q: %w", t.Name, err)
		}
		css := absURLs(string(src), filepath.Dir(file))
		for _, re := range []*regexp.Regexp{importRe, fontFaceRe} {
			css = re.ReplaceAllStringFunc(css, func(m string) string {
				top.WriteString(strings.TrimSpace(m) + "\n")
				return ""
			})
		}
		fmt.Fprintf(&scoped, "@scope (%s) {\n%s\n}\n", scope, css)
		return nil
	}
	if err := add(filepath.Join(t.Dir, "style.css"), "."+t.Class()); err != nil {
		return Stylesheet{}, err
	}
	for _, th := range themes {
		if err := add(filepath.Join(t.Dir, "themes", th+".css"), "."+t.Class()+"."+ThemeClass(th)); err != nil {
			return Stylesheet{}, err
		}
	}
	return Stylesheet{Top: top.String(), Scoped: scoped.String()}, nil
}

// absURLs resolves relative url()s against dir.
func absURLs(css, dir string) string {
	return urlRe.ReplaceAllStringFunc(css, func(m string) string {
		sub := urlRe.FindStringSubmatch(m)
		u := sub[2]
		if strings.Contains(u, ":") || strings.HasPrefix(u, "#") || strings.HasPrefix(u, "/") {
			return m
		}
		return `url("` + FileURL(filepath.Join(dir, filepath.FromSlash(u))) + `")`
	})
}
