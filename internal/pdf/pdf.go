// Package pdf prints an HTML file to PDF with a Chromium-based browser that
// is already installed (Chrome, Chromium, Edge, Brave). inkpress does not
// bundle a browser.
package pdf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/chromedp/cdproto/page"
	cdpruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// BrowserEnv names a browser executable, overriding the search.
const BrowserEnv = "INKPRESS_BROWSER"

// FindBrowser returns the path of an installed Chromium-based browser.
func FindBrowser() (string, error) {
	if p := os.Getenv(BrowserEnv); p != "" {
		return p, nil
	}
	var paths []string
	switch runtime.GOOS {
	case "darwin":
		for _, app := range []string{"Google Chrome", "Chromium", "Microsoft Edge", "Brave Browser", "Vivaldi", "Google Chrome Canary"} {
			paths = append(paths, filepath.Join("/Applications", app+".app", "Contents", "MacOS", app))
		}
	case "windows":
		for _, root := range []string{os.Getenv("ProgramFiles"), os.Getenv("ProgramFiles(x86)"), os.Getenv("LocalAppData")} {
			if root == "" {
				continue
			}
			paths = append(paths,
				filepath.Join(root, `Google\Chrome\Application\chrome.exe`),
				filepath.Join(root, `Microsoft\Edge\Application\msedge.exe`),
				filepath.Join(root, `BraveSoftware\Brave-Browser\Application\brave.exe`),
				filepath.Join(root, `Chromium\Application\chrome.exe`))
		}
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	for _, name := range []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser", "microsoft-edge", "microsoft-edge-stable", "brave-browser", "chrome"} {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("no Chromium-based browser found (Chrome, Chromium, Edge or Brave); install one, or set %s or --browser", BrowserEnv)
}

// Options configure one print.
type Options struct {
	// Browser is the executable; empty searches with FindBrowser.
	Browser string
	// Width and Height are the page size in CSS pixels.
	Width, Height float64
	Timeout       time.Duration
}

// Print loads the HTML file at path and returns it printed to PDF.
func Print(ctx context.Context, path string, opts Options) ([]byte, error) {
	browser := opts.Browser
	if browser == "" {
		var err error
		if browser, err = FindBrowser(); err != nil {
			return nil, err
		}
	}
	if opts.Timeout == 0 {
		opts.Timeout = 2 * time.Minute
	}
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(browser),
		// The page loads the pack's fonts and images from disk.
		chromedp.Flag("allow-file-access-from-files", true),
	)
	ctx, cancelAlloc := chromedp.NewExecAllocator(ctx, allocOpts...)
	defer cancelAlloc()
	// Errors come back from Run; chromedp's own log is noise.
	quiet := func(string, ...any) {}
	ctx, cancelTab := chromedp.NewContext(ctx, chromedp.WithErrorf(quiet), chromedp.WithLogf(quiet))
	defer cancelTab()

	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	u := "file://" + filepath.ToSlash(abs)
	if runtime.GOOS == "windows" {
		u = "file:///" + filepath.ToSlash(abs)
	}

	const pxPerInch = 96
	var buf []byte
	var ready bool
	err = chromedp.Run(ctx,
		chromedp.Navigate(u),
		// Wait for web fonts; images are in by the load event.
		chromedp.Evaluate(`document.fonts.ready.then(() => true)`, &ready,
			func(p *cdpruntime.EvaluateParams) *cdpruntime.EvaluateParams { return p.WithAwaitPromise(true) }),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithPaperWidth(opts.Width / pxPerInch).
				WithPaperHeight(opts.Height / pxPerInch).
				WithMarginTop(0).WithMarginBottom(0).WithMarginLeft(0).WithMarginRight(0).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("print with %s: %w", browser, err)
	}
	return buf, nil
}
