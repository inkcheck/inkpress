// Package deck loads a presentation pack: settings.toml and the slides.
package deck

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

// Settings is a pack's settings.toml.
type Settings struct {
	Title    string `toml:"title"`
	Author   string `toml:"author"`
	Date     string `toml:"date"`
	Template string `toml:"template"`
	Theme    string `toml:"theme"`
	// Output is the default PDF path, relative to the pack.
	Output string         `toml:"output"`
	Page   Page           `toml:"page"`
	Meta   map[string]any `toml:"meta"`
}

// Page is the slide size. Width and Height are CSS pixels (96 per inch).
// Set Width and Aspect, or Width and Height.
type Page struct {
	Aspect string  `toml:"aspect"`
	Width  float64 `toml:"width"`
	Height float64 `toml:"height"`
}

// Defaults for a pack whose settings leave them out.
const (
	DefaultTemplate = "default"
	DefaultAspect   = "16:9"
	DefaultWidth    = 1280
)

// ParseSettings reads settings.toml. Unknown keys are errors, so typos do
// not pass silently.
func ParseSettings(text string) (Settings, error) {
	var s Settings
	md, err := toml.Decode(text, &s)
	if err != nil {
		return s, err
	}
	if undec := md.Undecoded(); len(undec) > 0 {
		var keys []string
		for _, k := range undec {
			// Anything goes under [meta].
			if len(k) > 0 && k[0] != "meta" {
				keys = append(keys, k.String())
			}
		}
		if len(keys) > 0 {
			return s, fmt.Errorf("unknown settings: %s", strings.Join(keys, ", "))
		}
	}
	if s.Template == "" {
		s.Template = DefaultTemplate
	}
	if _, _, err := s.Page.Size(); err != nil {
		return s, err
	}
	return s, nil
}

// Size returns the slide width and height in CSS pixels.
func (p Page) Size() (w, h float64, err error) {
	w = p.Width
	if w == 0 {
		w = DefaultWidth
	}
	if w < 0 || p.Height < 0 {
		return 0, 0, fmt.Errorf("page: width and height must be positive")
	}
	if p.Height > 0 {
		if p.Aspect != "" {
			return 0, 0, fmt.Errorf("page: set aspect or height, not both")
		}
		return w, p.Height, nil
	}
	aspect := p.Aspect
	if aspect == "" {
		aspect = DefaultAspect
	}
	x, y, ok := strings.Cut(aspect, ":")
	ax, err1 := strconv.ParseFloat(strings.TrimSpace(x), 64)
	ay, err2 := strconv.ParseFloat(strings.TrimSpace(y), 64)
	if !ok || err1 != nil || err2 != nil || ax <= 0 || ay <= 0 {
		return 0, 0, fmt.Errorf("page: aspect %q is not W:H, e.g. 16:9", aspect)
	}
	return w, w * ay / ax, nil
}
