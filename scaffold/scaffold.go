// Package scaffold writes a new presentation pack: a sample deck and the
// default template.
package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed all:pack
var pack embed.FS

// New writes the sample pack into dir. With templateOnly it writes only
// .templates/default, to add the template to an existing pack. It never
// overwrites a file.
func New(dir string, templateOnly bool) ([]string, error) {
	root, _ := fs.Sub(pack, "pack")
	var files []string
	err := fs.WalkDir(root, ".", func(p string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() {
			return err
		}
		if templateOnly && !strings.HasPrefix(p, ".templates/") {
			return nil
		}
		files = append(files, p)
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			return nil, fmt.Errorf("%s already exists", filepath.Join(dir, f))
		}
	}
	for _, f := range files {
		data, err := fs.ReadFile(root, f)
		if err != nil {
			return nil, err
		}
		out := filepath.Join(dir, filepath.FromSlash(f))
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(out, data, 0o644); err != nil {
			return nil, err
		}
	}
	return files, nil
}
