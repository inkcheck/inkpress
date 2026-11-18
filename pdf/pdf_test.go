package pdf

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPrint(t *testing.T) {
	if testing.Short() {
		t.Skip("starts a browser")
	}
	if _, err := FindBrowser(); err != nil {
		t.Skip(err)
	}
	path := filepath.Join(t.TempDir(), "x.html")
	os.WriteFile(path, []byte(`<style>@page{size:400px 300px;margin:0}</style><p>one</p><p style="break-before:page">two</p>`), 0o644)
	buf, err := Print(context.Background(), path, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(buf, []byte("%PDF")) {
		t.Fatal("not a PDF")
	}
}
