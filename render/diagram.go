package render

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/inkcheck/inkpress/deck"
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
	var themeArgs []string
	// Match a light or dark slide, unless the diagram picks its own theme.
	if (theme == "light" || theme == "dark") && !strings.Contains(src, "theme:") {
		themeArgs = []string{"--theme", theme}
	}
	key := cacheKey(bin, s.Path, src, themeArgs)
	if fig, ok := r.cache.get(key, r.gen); ok {
		return fig, nil
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

	args := append(append([]string{"render", "-o", out}, themeArgs...), in)
	if msg, err := exec.Command(bin, args...).CombinedOutput(); err != nil {
		return "", fmt.Errorf("%s: diagram: %s", s.Path, strings.TrimSpace(strings.ReplaceAll(string(msg), in, "block")))
	}
	svg, err := os.ReadFile(out)
	if err != nil {
		return "", err
	}
	// An <img>, not inline SVG, keeps each diagram's IDs and styles apart.
	fig := fmt.Sprintf(`<figure class="diagram"><img src="data:image/svg+xml;base64,%s" alt=""></figure>`,
		base64.StdEncoding.EncodeToString(svg))
	r.cache.put(key, fig, r.gen)
	return fig, nil
}

// DiagramCache keeps rendered diagrams between renders of a deck. Once no
// render is running, it drops the diagrams the newest render did not use, so
// it holds about one deck's diagrams. It is safe for concurrent use. Files a
// diagram refers to, such as icons, are not part of the key: call Reset when
// they change.
type DiagramCache struct {
	mu      sync.Mutex
	figures map[string]cached
	gen     int // the newest render
	running int
}

type cached struct {
	figure string
	gen    int // the newest render that used it
}

// NewDiagramCache returns an empty cache.
func NewDiagramCache() *DiagramCache {
	return &DiagramCache{figures: map[string]cached{}}
}

// Reset empties the cache.
func (c *DiagramCache) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	clear(c.figures)
}

// Len is the number of diagrams in the cache.
func (c *DiagramCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.figures)
}

func cacheKey(bin, slide, src string, args []string) string {
	h := sha256.New()
	// The slide's directory matters, since icons resolve against it.
	for _, part := range append([]string{bin, filepath.Dir(slide), src}, args...) {
		h.Write([]byte(part))
		h.Write([]byte{0})
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// begin starts a render and returns its generation.
func (c *DiagramCache) begin() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.gen++
	c.running++
	return c.gen
}

func (c *DiagramCache) end() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.running--
	if c.running > 0 {
		return
	}
	for k, f := range c.figures {
		if f.gen < c.gen {
			delete(c.figures, k)
		}
	}
}

func (c *DiagramCache) get(key string, gen int) (string, bool) {
	if c == nil {
		return "", false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	f, ok := c.figures[key]
	if ok {
		f.gen = max(f.gen, gen)
		c.figures[key] = f
	}
	return f.figure, ok
}

func (c *DiagramCache) put(key, fig string, gen int) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.figures[key] = cached{figure: fig, gen: max(c.figures[key].gen, gen)}
}
