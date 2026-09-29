package render

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/inkcheck/inkpress/internal/deck"
)

// Diagrams come from inkline (https://github.com/inkcheck/inkline), run as a
// command: a ```inkline block, or a ```d2 block with an inkline header,
// becomes an SVG image. Without inkline on the PATH the block stays a code
// block. ```verf blocks and verf headers, from before inkline was renamed,
// still work.

// InklineEnv names the inkline executable, overriding the PATH.
const InklineEnv = "INKPRESS_INKLINE"

type inkline struct {
	path    string
	looked  bool
	missing bool
	tmp     string
}

func (v *inkline) find() (string, bool) {
	if !v.looked {
		v.looked = true
		v.path = os.Getenv(InklineEnv)
		if v.path == "" {
			v.path, _ = exec.LookPath("inkline")
		}
		v.missing = v.path == ""
	}
	return v.path, !v.missing
}

func (v *inkline) cleanup() {
	if v.tmp != "" {
		os.RemoveAll(v.tmp)
	}
}

// diagrams replaces diagram blocks in the slide body with images.
func (r *renderer) diagrams(s *deck.Slide, theme string) (string, error) {
	lines := strings.Split(s.Body, "\n")
	var out []string
	for i := 0; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i])
		fence := deck.FenceOpen(t)
		if fence == "" {
			out = append(out, lines[i])
			continue
		}
		lang, _, _ := strings.Cut(strings.TrimSpace(t[len(fence):]), " ")
		end := i + 1
		for end < len(lines) {
			c := strings.TrimSpace(lines[end])
			if strings.HasPrefix(c, fence) && strings.Trim(c, fence[:1]) == "" {
				break
			}
			end++
		}
		src := strings.Join(lines[i+1:min(end, len(lines))], "\n")
		isDiagram := lang == "inkline" || lang == "verf" ||
			(lang == "d2" && (strings.Contains(src, "inkline:") || strings.Contains(src, "verf:")))
		if !isDiagram {
			out = append(out, lines[i:min(end+1, len(lines))]...)
			i = end
			continue
		}
		img, err := r.diagram(s, src, theme)
		if err != nil {
			return "", err
		}
		if img == "" {
			// No inkline: show the source.
			if lang == "inkline" || lang == "verf" {
				lines[i] = strings.Replace(lines[i], lang, "d2", 1)
			}
			out = append(out, lines[i:min(end+1, len(lines))]...)
		} else {
			out = append(out, "", img, "")
		}
		i = end
	}
	return strings.Join(out, "\n"), nil
}

// diagram renders D2 source with inkline and returns an HTML figure, or ""
// when inkline is not installed.
func (r *renderer) diagram(s *deck.Slide, src, theme string) (string, error) {
	bin, ok := r.inkline.find()
	if !ok {
		r.warn("%s: inkline not found, showing the diagram source; install it from https://github.com/inkcheck/inkline", s.Path)
		return "", nil
	}
	// A verf header becomes an inkline one.
	if !strings.Contains(src, "inkline:") {
		src = strings.Replace(src, "verf:", "inkline:", 1)
	}
	// Icons in the diagram resolve against the D2 file's directory, so the
	// file goes next to the slide.
	sum := sha256.Sum256([]byte(src))
	in := filepath.Join(r.deck.Dir, filepath.Dir(s.Path), fmt.Sprintf(".inkpress-%x.d2", sum[:6]))
	if err := os.WriteFile(in, []byte(src), 0o644); err != nil {
		return "", err
	}
	defer os.Remove(in)
	if r.inkline.tmp == "" {
		dir, err := os.MkdirTemp("", "inkpress-inkline-")
		if err != nil {
			return "", err
		}
		r.inkline.tmp = dir
	}
	out := filepath.Join(r.inkline.tmp, fmt.Sprintf("%x.svg", sum[:6]))

	args := []string{"render", "-o", out}
	// Match a light or dark slide, unless the diagram picks its own theme.
	if (theme == "light" || theme == "dark") && !strings.Contains(src, "theme:") {
		args = append(args, "--theme", theme)
	}
	args = append(args, in)
	if msg, err := exec.Command(bin, args...).CombinedOutput(); err != nil {
		return "", fmt.Errorf("%s: diagram: %s", s.Path, strings.TrimSpace(strings.ReplaceAll(string(msg), in, "block")))
	}
	svg, err := os.ReadFile(out)
	if err != nil {
		return "", err
	}
	// An <img>, not inline SVG, keeps each diagram's IDs and styles apart.
	return fmt.Sprintf(`<figure class="diagram"><img src="data:image/svg+xml;base64,%s" alt=""></figure>`,
		base64.StdEncoding.EncodeToString(svg)), nil
}
