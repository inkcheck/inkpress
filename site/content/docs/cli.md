---
title: CLI
weight: 8
---

```sh
inkpress new [--template-only] <dir>   # start a pack: sample slides and the default template
inkpress pdf [flags] [dir]             # build the PDF (default dir: .)
inkpress html [flags] [dir]            # build an HTML preview
inkpress list [dir]                    # list the slides, their layouts and themes
inkpress version                       # print the version
```

Flags and the directory can come in any order.

| Flag | Command | Meaning |
|---|---|---|
| `-o <path>` | `pdf`, `html` | output file; default: settings `output`, or `<dir name>.pdf` in the pack |
| `--browser <path>` | `pdf` | the browser to print with; default: search, or `$INKPRESS_BROWSER` |
| `--timeout <duration>` | `pdf` | give up printing after this long; default `2m` |
| `--template-only` | `new` | write only `.templates/default`, into an existing pack |

## list

```console
$ inkpress list talk
  1  01-title.md                  default/title/light  notes
  2  02-agenda.md                 default/agenda/light  notes
  3  03-why.md                    default/content/light  notes
  4  04-section.md                default/section/light
```

Each line shows the slide's number, file, template, layout and theme, and
whether it has notes.

## Errors and warnings

An error names the slide file, such as
`slides/04-x.md: template "default" has no layout "twocol" (layouts: …)`.

Warnings, such as a missing image or inkline not being installed, do not stop
the build.

## The browser

inkpress renders the deck to one HTML page. It then prints that page with a
Chromium-based browser that is already installed, over the DevTools protocol.
It looks for Chrome, Chromium, Edge, Brave and Vivaldi.

To choose a browser, use `--browser <path>` or set `INKPRESS_BROWSER`. If
inkpress finds none, it fails with "no Chromium-based browser found". Some
browsers, such as Brave Origin, turn printing off and fail with "Printing is
not available".

## Environment

| Variable | Meaning |
|---|---|
| `INKPRESS_BROWSER` | the browser to print with |
| `INKPRESS_INKLINE` | the inkline executable for diagrams; default: `inkline` on the PATH |
