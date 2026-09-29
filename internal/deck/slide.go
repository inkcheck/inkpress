package deck

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Slide is one file in slides/.
type Slide struct {
	// Path is the file's path on disk.
	Path string
	// Meta is the frontmatter. The keys inkpress reads are in Known;
	// templates can read the rest.
	Meta map[string]any
	// Body is the Markdown between the frontmatter and the notes.
	Body string
	// Notes is the tailmatter: speaker notes, left out of the presentation.
	Notes string
}

// Frontmatter keys inkpress reads itself.
const (
	KeyTemplate   = "template"
	KeyLayout     = "layout"
	KeyTheme      = "theme"
	KeyTitle      = "title"
	KeyClass      = "class"
	KeyBackground = "background"
	KeyImage      = "image"
	KeyHeader     = "header"
	KeyFooter     = "footer"
	KeySkip       = "skip"
)

// NotesMarker starts the tailmatter: every line after it is a note.
const NotesMarker = "!--"

// ColumnMarker splits a slide body into columns.
const ColumnMarker = "|||"

// ParseSlide splits a slide file into frontmatter, body and notes.
// Frontmatter is YAML between --- lines, or TOML between +++ lines.
func ParseSlide(path, text string) (*Slide, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	s := &Slide{Path: path, Meta: map[string]any{}}

	lines := strings.Split(text, "\n")
	if len(lines) > 0 && (lines[0] == "---" || lines[0] == "+++") {
		fence := lines[0]
		end := -1
		for i := 1; i < len(lines); i++ {
			if strings.TrimRight(lines[i], " \t") == fence {
				end = i
				break
			}
		}
		if end < 0 {
			return nil, fmt.Errorf("%s: frontmatter has no closing %s", path, fence)
		}
		fm := strings.Join(lines[1:end], "\n")
		var err error
		if fence == "---" {
			err = yaml.Unmarshal([]byte(fm), &s.Meta)
		} else {
			_, err = toml.Decode(fm, &s.Meta)
		}
		if err != nil {
			return nil, fmt.Errorf("%s: frontmatter: %w", path, err)
		}
		if s.Meta == nil {
			s.Meta = map[string]any{}
		}
		lines = lines[end+1:]
	}

	body := lines
	for _, i := range Markers(lines, NotesMarker) {
		body = lines[:i]
		s.Notes = strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
		break
	}
	s.Body = strings.TrimSpace(strings.Join(body, "\n"))
	return s, nil
}

// String returns the frontmatter value for key, or "".
func (s *Slide) String(key string) string {
	switch v := s.Meta[key].(type) {
	case nil:
		return ""
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}

// Bool returns the frontmatter value for key, or def when it is not set.
func (s *Slide) Bool(key string, def bool) bool {
	if v, ok := s.Meta[key].(bool); ok {
		return v
	}
	return def
}

// Markers returns the indexes of lines equal to marker, outside fenced code
// blocks.
func Markers(lines []string, marker string) []int {
	var idx []int
	fence := ""
	for i, l := range lines {
		t := strings.TrimSpace(l)
		if fence != "" {
			if strings.HasPrefix(t, fence) && strings.Trim(t, fence[:1]) == "" {
				fence = ""
			}
			continue
		}
		if f := FenceOpen(t); f != "" {
			fence = f
			continue
		}
		if t == marker {
			idx = append(idx, i)
		}
	}
	return idx
}

// FenceOpen returns the fence (``` or ~~~, three or more) that opens a
// fenced code block on line t, or "".
func FenceOpen(t string) string {
	for _, c := range []string{"`", "~"} {
		n := len(t) - len(strings.TrimLeft(t, c))
		if n >= 3 {
			return t[:n]
		}
	}
	return ""
}
