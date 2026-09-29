package deck

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Deck is a loaded presentation pack.
type Deck struct {
	// Dir is the pack's absolute path.
	Dir      string
	Settings Settings
	// Slides are in filename order, without skipped slides.
	Slides []*Slide
}

// Pack layout.
const (
	SettingsFile = "settings.toml"
	SlidesDir    = "slides"
	AssetsDir    = "assets"
	TemplatesDir = ".templates"
)

// Load reads the pack in dir.
func Load(dir string) (*Deck, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	d := &Deck{Dir: abs}

	text, err := os.ReadFile(filepath.Join(abs, SettingsFile))
	switch {
	case os.IsNotExist(err):
		// Every setting has a default.
		d.Settings, err = ParseSettings("")
	case err == nil:
		d.Settings, err = ParseSettings(string(text))
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", SettingsFile, err)
	}

	entries, err := os.ReadDir(filepath.Join(abs, SlidesDir))
	if err != nil {
		return nil, fmt.Errorf("%s is not a presentation pack: %w", dir, err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".md") && !strings.HasPrefix(e.Name(), ".") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		path := filepath.Join(abs, SlidesDir, name)
		text, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		s, err := ParseSlide(filepath.Join(SlidesDir, name), string(text))
		if err != nil {
			return nil, err
		}
		if !s.Bool(KeySkip, false) {
			d.Slides = append(d.Slides, s)
		}
	}
	if len(d.Slides) == 0 {
		return nil, fmt.Errorf("%s: no slides", filepath.Join(dir, SlidesDir))
	}
	return d, nil
}

// Title is the deck title, or the pack directory's name.
func (d *Deck) Title() string {
	if d.Settings.Title != "" {
		return d.Settings.Title
	}
	return filepath.Base(d.Dir)
}
