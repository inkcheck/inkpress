# inkpress

inkpress turns a folder of Markdown slides into a PDF. Each slide is one
Markdown file. A template pack of HTML layouts and CSS themes controls the
design, and a Chromium-based browser that is already installed prints the
result.

```sh
brew install --cask inkcheck/tap/inkpress                   # macOS and Linux
go install github.com/inkcheck/inkpress/cmd/inkpress@latest    # or from source

inkpress new talk          # a sample pack with the default template
inkpress pdf talk          # -> talk/talk.pdf
inkpress html talk         # -> talk/talk.html, a preview for the browser
inkpress list talk         # slides in order, with template/layout/theme
```

## A presentation pack

```
talk/
  settings.toml           title, author, template, theme, page size
  README.md
  slides/                 one slide per file, in filename order: 01-intro.md, 02-…
  assets/                 images the slides use (jpg, png, svg)
  .templates/default/     a template pack; a deck can hold several
```

### settings.toml

```toml
title = "Quarterly review"
author = "A. Presenter"
date = "2026-09-29"
template = "default"      # a directory in .templates/
theme = "light"           # a theme of that template
output = "review.pdf"     # default: <directory name>.pdf

[page]
width = 1280              # CSS pixels, 96 per inch
aspect = "16:9"           # 16:9, 16:10, 4:3, any W:H; or set height instead

[meta]                    # anything; templates read it as .Deck.Meta
```

Every key is optional. An unknown key is an error, so typos are caught.

### A slide

```markdown
---
layout: two-column
theme: dark
---

# Heading

Left column

|||

Right column

!--
Speaker notes. Everything after the !-- line stays out of the presentation.
```

- **Frontmatter** is YAML between `---` lines, or TOML between `+++` lines.
  inkpress reads these keys:

  | Key | Meaning |
  |---|---|
  | `template` | the template pack, overriding settings |
  | `layout` | a layout of the template; default from `template.toml` |
  | `theme` | a theme of the template |
  | `title` | the slide title; default: the leading heading |
  | `class` | extra CSS classes on the slide |
  | `background` | a colour, gradient, or image path |
  | `image` | an image path, for layouts that use one |
  | `header`, `footer` | `false` hides them |
  | `skip` | `true` leaves the slide out |

  Layouts can read any other key as `.Meta.<key>`, for example `subtitle`,
  `kicker`, `cite`, `role` or `current` in the default template.
- **Tailmatter** is everything after a line that holds only `!--`. These are
  speaker notes, and they are not rendered.
- **Columns** are separated by lines that hold only `|||`.
- **Boxes**: Markdown between `:::` lines goes in a div. `::: box note` opens
  `<div class="box note">`, and a bare `:::` closes it. Boxes nest, and the
  template styles the classes. The default template has `box`, `card`,
  `row` (blocks side by side), `stat` (a big number), `figure` (an image and
  its caption), `gallery` (images in a row), `lead`, `muted` and `center`,
  and the tones `note`, `tip`, `important`, `warning` and `caution`. A box
  cannot span a `|||` line; one left open closes at the end of its column.
- **Callouts** use GitHub's syntax: a blockquote whose first line is
  `[!NOTE]`, `[!TIP]`, `[!IMPORTANT]`, `[!WARNING]` or `[!CAUTION]` becomes
  `<blockquote class="callout note">` with a `.callout-title` paragraph.
  Text after the marker, as in `> [!TIP] Try this`, replaces the title.
- **Markdown** is GitHub-flavoured: tables, task lists, strikethrough and
  autolinks. It also allows raw HTML. Code blocks are highlighted.
- **Images** resolve against the slide (`../assets/chart.svg`), then the pack
  (`assets/chart.svg`).
- **Diagrams**: a ` ```inkline ` block, or a ` ```d2 ` block with an inkline
  header, is rendered to SVG by [inkline](https://github.com/inkcheck/inkline).
  Light and dark slides get light and dark diagrams, unless the header sets a
  theme. Without inkline on the PATH, the slide shows the source and inkpress
  warns you. Set `INKPRESS_INKLINE` to use a particular inkline executable.
  ` ```verf ` blocks and verf headers from older decks still work.

## Template packs

```
.templates/<name>/
  template.toml           description, default layout and theme, code styles
  style.css               styles for .slide and everything inside it
  themes/<theme>.css      custom properties for .slide, one file per theme
  layouts/<layout>.html   Go html/template, one per layout
  partials/header.html    optional; rendered into .Header
  partials/footer.html    optional; rendered into .Footer
  assets/, fonts/         anything the CSS and layouts refer to
```

The default template has these layouts: `title`, `section`, `agenda`,
`content`, `two-column`, `grid`, `statement`, `image`, `image-left`,
`image-full`, `quote` and `blank`. It has two themes, `light` and
`dark`. To add the default template to an existing pack, run
`inkpress new --template-only`.

Layouts and partials get these fields:

| Field | |
|---|---|
| `.Heading` | the slide's leading `#` or `##` heading, as HTML |
| `.Body` | the rest of the slide |
| `.Content` | `.Heading` and `.Body` together |
| `.Columns` | `.Body` split at `\|\|\|` lines |
| `.Title`, `.Number`, `.Total` | the slide title, its number from 1, and the slide count |
| `.Layout`, `.Theme`, `.Template` | the names in use |
| `.Image` | the frontmatter `image`, resolved |
| `.Meta` | the slide's frontmatter |
| `.Header`, `.Footer` | the rendered partials; empty when the slide hides them |
| `.Deck.Title`, `.Deck.Author`, `.Deck.Date`, `.Deck.Meta` | from settings.toml |

The function `asset "logo.svg"` returns the URL of a file in the template's
`assets/`.

A deck can mix templates and themes, so inkpress wraps each stylesheet in
`@scope (.tpl-<name>)` and each theme in `@scope (.tpl-<name>.theme-<theme>)`.
It moves `@font-face` and `@import` to the top level, and makes relative
`url()`s absolute. Each slide sits in a `.page` that is a CSS size container,
so size things in `cqh` and `cqw` (1% of the slide's height or width). The
design then scales with the page size.

## PDF

inkpress renders the deck to one HTML page. It then prints that page with a
Chromium-based browser that is already installed, over the DevTools protocol.
It looks for Chrome, Chromium, Edge, Brave and Vivaldi. To choose a browser,
use `--browser <path>` or `INKPRESS_BROWSER`. inkpress does not bundle a
browser.

## Layout

```
cmd/inkpress/         CLI
internal/deck/        settings, slide parsing (frontmatter, notes, columns), loading
internal/templates/   template packs, CSS scoping
internal/render/      Markdown to HTML, layouts, inkline diagrams
internal/pdf/         browser discovery and printing
internal/scaffold/    the pack `inkpress new` writes, with the default template
skills/inkpress/      Agent skill: teaches Claude to write decks and templates
scripts/              notices, releases
```

```sh
make build            # bin/inkpress
make test
make demo             # examples/demo and its PDF
```

## Agent skill

`skills/inkpress/` is an [Agent Skill](https://docs.claude.com/en/docs/agents-and-tools/agent-skills).
It teaches Claude to plan a deck, write the slides with the right layouts and
notes, build the PDF, and change templates. To use it in Claude Code, copy or
symlink it into `~/.claude/skills/` for all projects, or into a project's
`.claude/skills/`:

```sh
ln -s "$PWD/skills/inkpress" ~/.claude/skills/inkpress
```

It works with the [inkline skill](https://github.com/inkcheck/inkline) for diagrams.
