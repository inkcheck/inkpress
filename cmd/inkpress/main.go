// Command inkpress turns a pack of Markdown slides into a PDF.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/inkcheck/inkpress/internal/deck"
	"github.com/inkcheck/inkpress/internal/pdf"
	"github.com/inkcheck/inkpress/internal/render"
	"github.com/inkcheck/inkpress/internal/scaffold"
	"github.com/inkcheck/inkpress/internal/templates"
)

const usageText = `inkpress turns Markdown slides into a PDF.

Usage:
  inkpress new [--template-only] <dir>   start a pack: sample slides and the default template
  inkpress pdf [flags] [dir]             build the PDF (default dir: .)
  inkpress html [flags] [dir]            build an HTML preview
  inkpress list [dir]                    list the slides, their layouts and themes
  inkpress version                       print the version

A pack:

  settings.toml        title, author, template, theme, page size
  slides/*.md          one slide per file, in filename order; frontmatter
                       picks layout and theme; notes follow a !-- line
  assets/              images (jpg, png, svg)
  .templates/<name>/   layouts, themes, header, footer, fonts

PDFs are printed by an installed Chromium-based browser: Chrome, Chromium,
Edge or Brave. Set --browser or INKPRESS_BROWSER to choose one.
Diagrams in ` + "```inkline" + ` blocks need inkline: https://github.com/inkcheck/inkline.

Flags:
`

// Set by the release build (see .goreleaser.yaml).
var (
	version = ""
	commit  = ""
	date    = ""
)

func versionString() string {
	v := version
	if v == "" {
		v = "dev"
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			v = strings.TrimPrefix(info.Main.Version, "v")
		}
	}
	if commit != "" {
		v += fmt.Sprintf(" (%s, %s)", commit, date)
	}
	return v
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "inkpress:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, usageText)
		return fmt.Errorf("missing command")
	}
	cmd, args := args[0], args[1:]
	switch cmd {
	case "help", "-h", "--help":
		fmt.Print(usageText)
		return nil
	case "version", "--version", "-v":
		fmt.Println("inkpress", versionString())
		return nil
	case "new", "pdf", "html", "list":
	default:
		fmt.Fprint(os.Stderr, usageText)
		return fmt.Errorf("unknown command %q", cmd)
	}

	fl := flag.NewFlagSet(cmd, flag.ContinueOnError)
	fl.Usage = func() { fmt.Fprint(os.Stderr, usageText); fl.PrintDefaults() }
	var out, browser *string
	var timeout *time.Duration
	var templateOnly *bool
	switch cmd {
	case "pdf":
		out = fl.String("o", "", "output PDF (default: settings output, or <dir name>.pdf in the pack)")
		browser = fl.String("browser", "", "Chromium-based browser executable (default: search, or $"+pdf.BrowserEnv+")")
		timeout = fl.Duration("timeout", 2*time.Minute, "give up printing after this long")
	case "html":
		out = fl.String("o", "", "output HTML (default: <dir name>.html in the pack)")
	case "new":
		templateOnly = fl.Bool("template-only", false, "write only .templates/default, into an existing pack")
	}
	files, err := parseArgs(fl, args)
	if err != nil || files == nil {
		return err
	}
	if len(files) > 1 {
		fl.Usage()
		return fmt.Errorf("expected one directory")
	}
	dir := "."
	if len(files) == 1 {
		dir = files[0]
	}

	switch cmd {
	case "new":
		if len(files) == 0 && !*templateOnly {
			fl.Usage()
			return fmt.Errorf("new: name the directory for the pack")
		}
		written, err := scaffold.New(dir, *templateOnly)
		if err != nil {
			return err
		}
		fmt.Printf("wrote %d files to %s\n", len(written), dir)
		if !*templateOnly {
			fmt.Printf("next: inkpress pdf %s\n", dir)
		}
		return nil
	case "list":
		return list(dir)
	}

	d, err := deck.Load(dir)
	if err != nil {
		return err
	}
	res, err := render.HTML(d)
	if err != nil {
		return err
	}
	for _, w := range res.Warnings {
		fmt.Fprintln(os.Stderr, "inkpress: warning:", w)
	}
	base := filepath.Base(d.Dir)

	if cmd == "html" {
		path := *out
		if path == "" {
			path = filepath.Join(d.Dir, base+".html")
		}
		if err := os.WriteFile(path, res.HTML, 0o644); err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	}

	path := *out
	if path == "" {
		name := d.Settings.Output
		if name == "" {
			name = base + ".pdf"
		}
		path = filepath.Join(d.Dir, name)
	}
	// The browser loads the HTML from disk; a temporary file is enough, since
	// every URL in it is absolute.
	tmp, err := os.CreateTemp("", "inkpress-*.html")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(res.HTML); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	buf, err := pdf.Print(context.Background(), tmp.Name(), pdf.Options{
		Browser: *browser, Width: res.Width, Height: res.Height, Timeout: *timeout,
	})
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, buf, 0o644); err != nil {
		return err
	}
	fmt.Printf("%s (%d slides)\n", path, len(d.Slides))
	return nil
}

// parseArgs parses flags and arguments in any order, so both of these work:
//
//	inkpress pdf -o out.pdf deck
//	inkpress pdf deck -o out.pdf
//
// Everything after -- is an argument. It returns nil files, and no error,
// for -h.
func parseArgs(fl *flag.FlagSet, args []string) ([]string, error) {
	files := []string{}
	for rest := args; len(rest) > 0; {
		if err := fl.Parse(rest); errors.Is(err, flag.ErrHelp) {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		parsed := rest[:len(rest)-fl.NArg()]
		rest = fl.Args()
		if len(parsed) > 0 && parsed[len(parsed)-1] == "--" {
			files = append(files, rest...)
			break
		}
		if len(rest) > 0 {
			files, rest = append(files, rest[0]), rest[1:]
		}
	}
	return files, nil
}

// list prints each slide with its template, layout and theme.
func list(dir string) error {
	d, err := deck.Load(dir)
	if err != nil {
		return err
	}
	loaded := map[string]*templates.Template{}
	for i, s := range d.Slides {
		name := firstOf(s.String(deck.KeyTemplate), d.Settings.Template)
		t, ok := loaded[name]
		if !ok {
			if t, err = templates.Load(filepath.Join(d.Dir, deck.TemplatesDir), name); err != nil {
				return err
			}
			loaded[name] = t
		}
		layout := firstOf(s.String(deck.KeyLayout), t.Config.Layout)
		theme := firstOf(s.String(deck.KeyTheme), d.Settings.Theme, t.Config.Theme)
		notes := ""
		if s.Notes != "" {
			notes = "  notes"
		}
		fmt.Printf("%3d  %-28s %s/%s/%s%s\n", i+1, filepath.Base(s.Path), name, layout, theme, notes)
	}
	return nil
}

func firstOf(s ...string) string {
	for _, x := range s {
		if x != "" {
			return x
		}
	}
	return ""
}
