// Package inkdeck writes .inkdeck files: a whole pack in one zip, for
// sharing. The zip holds settings.toml, slides/, assets/ and .templates/,
// inside one folder named after the pack. InkStudio opens them, and unzip
// does too.
package inkdeck

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/inkcheck/inkpress/deck"
)

// Ext is the extension of shared decks.
const Ext = ".inkdeck"

// Skip reports whether a file of the pack stays out of a shared deck:
// hidden files other than the templates, links, and what builds and exports
// leave at the top of the pack. rel is the slash-separated path in the pack.
func Skip(rel string, d fs.DirEntry) bool {
	parts := strings.Split(rel, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ".") && !(i == 0 && part == deck.TemplatesDir) {
			return true
		}
	}
	if d.Type()&fs.ModeSymlink != 0 {
		return true
	}
	if len(parts) == 1 && !d.IsDir() {
		switch strings.ToLower(path.Ext(rel)) {
		case ".pdf", ".pptx", ".odp", ".html", Ext:
			return true
		}
	}
	return false
}

// Share zips the pack in dir into an .inkdeck at out. out may be at the top
// of the pack, which Skip leaves out, but not deeper in it. It writes a
// temporary file next to out and renames it, so a failure leaves nothing.
// It returns the number of files it packed.
func Share(dir, out string) (int, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return 0, err
	}
	if out, err = filepath.Abs(out); err != nil {
		return 0, err
	}
	if st, err := os.Stat(filepath.Join(dir, deck.SlidesDir)); err != nil || !st.IsDir() {
		return 0, errors.New("no slides/ folder in " + dir)
	}
	if rel, err := filepath.Rel(dir, filepath.Dir(out)); err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return 0, errors.New("save the .inkdeck outside the pack's folders")
	}

	root := filepath.Base(dir)
	tmp, err := os.CreateTemp(filepath.Dir(out), ".inkpress-share-*")
	if err != nil {
		return 0, err
	}
	defer os.Remove(tmp.Name())
	zw := zip.NewWriter(tmp)
	count := 0
	err = filepath.WalkDir(dir, func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if file == dir {
			return nil
		}
		rel, err := filepath.Rel(dir, file)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if Skip(rel, d) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() || !d.Type().IsRegular() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		h.Name = root + "/" + rel
		h.Method = zip.Deflate
		w, err := zw.CreateHeader(h)
		if err != nil {
			return err
		}
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(w, f)
		count++
		return err
	})
	if cerr := zw.Close(); err == nil {
		err = cerr
	}
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return 0, err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return 0, err
	}
	return count, os.Rename(tmp.Name(), out)
}
