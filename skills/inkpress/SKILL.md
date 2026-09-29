---
name: inkpress
description: Write presentations as Markdown and build them into a PDF with inkpress. An inkpress pack is a folder with settings.toml, one Markdown file per slide in slides/, images in assets/ and template packs (HTML layouts, CSS themes, fonts) in .templates/. Use this skill whenever the user wants a slide deck, presentation, talk, pitch or keynote, wants to turn notes, a document or an outline into slides, edits files in a folder with settings.toml and slides/*.md, asks for a new slide layout, theme or template, or mentions inkpress at all. Use it even if they only say "make some slides", "deck", "export the talk to PDF" or "speaker notes".
---

# inkpress

inkpress turns a folder of Markdown slides into a PDF. Each slide is one
Markdown file. Its frontmatter picks a layout and theme from a template pack,
and speaker notes follow a `!--` line. inkpress renders the deck to HTML, and
an installed Chromium-based browser (Chrome, Chromium, Edge or Brave) prints
it.

## Workflow

1. **Start a pack** with `inkpress new <dir>`, unless one exists. This writes
   sample slides and the default template to `.templates/default`. Replace the
   sample slides with the user's content. Do not hand-write a template from
   scratch unless the user asks for a new design.
2. **Plan the deck** before writing: one idea per slide. Use a title slide,
   section slides between parts, and a closing slide. A 10-minute talk has
   about 8–12 slides.
3. **Write one file per slide** in `slides/`, named with a sequence number:
   `01-title.md`, `02-problem.md`, `03-…`. Slides run in filename order. Leave
   gaps in the numbers (`10-`, `20-`) if the user expects to insert slides.
4. **Put speaker notes after `!--`** on every content slide when the user
   gives more detail than fits on the slide. Keep the slide terse and move the
   rest into the notes.
5. **Build and check**: `inkpress list <dir>`, then `inkpress pdf <dir>`. If you
   can view images, look at the pages (for example with
   `pdftoppm -r 50 -png deck.pdf page`) and fix any slide whose text
   overflows. The slide does not scroll; overflowing text is cut off.

## A pack

```
deck/
  settings.toml          title, author, date, template, theme, page size
  README.md
  slides/NN-name.md      one slide per file
  assets/                images: jpg, png, svg
  .templates/<name>/     template packs
```

```toml
# settings.toml. Every key is optional; an unknown key is an error.
title = "Quarterly review"
author = "A. Presenter"
date = "2026-09-29"
template = "default"
theme = "light"            # light or dark in the default template
output = "review.pdf"      # default: <directory name>.pdf

[page]
width = 1280               # CSS px; 96 per inch
aspect = "16:9"            # or 4:3, 16:10, any W:H; or height = … instead

[meta]                     # free-form; templates read .Deck.Meta
```

## A slide

```markdown
---
layout: two-column
---

# The heading

Left column

|||

Right column

!--
Speaker notes: everything after the !-- line. Never rendered.
```

- **Frontmatter**: YAML between `---` lines (or TOML between `+++` lines).

  | Key | Meaning |
  |---|---|
  | `layout` | a layout of the slide's template; default `content` |
  | `theme` | a theme of the template, overriding settings |
  | `template` | another template pack in `.templates/` |
  | `title` | slide title; default is the leading heading |
  | `class` | extra CSS classes on the slide |
  | `background` | a colour, a gradient, or an image path (cover) |
  | `image` | image path for layouts that use one (`image`) |
  | `header`, `footer` | `false` hides them |
  | `skip` | `true` leaves the slide out of the build |

  Layouts can read any other key as `.Meta.<key>`.
- **The leading `#` or `##` heading** becomes `.Heading`. Layouts place it
  separately from the body, so start each content slide with one heading.
- **`|||` lines** split the body into columns.
- **`!--` line** starts the notes. It is ignored inside code blocks.
- **Markdown** is GitHub-flavoured: tables, task lists, strikethrough,
  autolinks and raw HTML. Code blocks with a language are highlighted.
- **Images** go in `assets/`. Refer to them as `../assets/x.png` (relative to
  the slide, so editor previews work) or `assets/x.png`.

## Default template layouts

| Layout | For | Reads |
|---|---|---|
| `title` | the opening slide | `# Title`, `subtitle`; shows author and date from settings; no header or footer |
| `section` | a divider between parts, on an accent background | `# Section name`, `kicker` (small text above) |
| `content` | the default: heading and body | |
| `two-column` | comparisons, text beside code or an image | columns split by `\|\|\|` |
| `image` | text on the left half, a picture filling the right half | `image: ../assets/x.jpg` |
| `quote` | one large quotation | a `>` blockquote, `cite` |
| `blank` | full-bleed content, no padding, no header or footer | |

Themes: `light` and `dark`. Set one per slide with `theme: dark`, for example
for a closing slide or a quote.

## Fitting content

The slide is 1280×720 at the default size, and the body text is about 24px.
As a guide, a `content` slide holds:

- a heading and up to 6 bullets of one line each, or
- a heading and about 60 words of prose, or
- a heading and a code block of about 12 lines, or a table of about 6 rows.

Split anything longer across slides, or move the detail into the notes.
Prefer bullets that are phrases, not sentences. Do not shrink text with
inline styles to make it fit.

## Diagrams

A ```` ```inkline ```` block, or a ```` ```d2 ```` block with an inkline
header, is rendered to SVG by [inkline](https://github.com/inkcheck/inkline)
and scaled to fit the slide. Use the `inkline` skill to write the diagram.
Leave `theme` out of the inkline header so that the diagram follows the
slide's light or dark theme.
Prefer the default top-to-bottom direction; a wide diagram in a row shrinks
to fit the width and becomes hard to read.

````markdown
# Checkout flow

```inkline
vars: { inkline: { diagram: flowchart; design_system: carbon } }
cart: Cart {class: start}
pay: Pay {class: process}
done: Done {class: end}
cart -> pay -> done
```
````

If inkline is not installed, inkpress warns and shows the D2 source instead.

## Commands

```sh
inkpress new <dir>                  # sample pack + default template
inkpress new --template-only [dir]  # add .templates/default to an existing pack
inkpress list [dir]                 # slides in order: template/layout/theme, notes
inkpress pdf [dir] [-o out.pdf] [--browser path] [--timeout 2m]
inkpress html [dir] [-o out.html]   # preview; open it in a browser
```

- Flags and the directory go in any order. The directory defaults to `.`.
- **Errors** name the slide file, such as
  `slides/04-x.md: template "default" has no layout "twocol" (layouts: …)`.
- **Warnings** (a missing image, inkline not installed) do not stop the build.
  Read them and fix them.
- **"no Chromium-based browser found"**: install Chrome, or pass
  `--browser <path>` or set `INKPRESS_BROWSER`. Some browsers, such as Brave
  Origin, turn printing off and fail with "Printing is not available".
- **If `inkpress` is not on the PATH**: `brew install --cask
  inkcheck/tap/inkpress`, `go install
  github.com/inkcheck/inkpress/cmd/inkpress@latest`, or `make build` in the
  inkpress repository for `bin/inkpress`.

## New layouts, themes and templates

When the user wants a design change, edit the pack's template in
`.templates/<name>/`: add a layout, adjust a theme's colours, or copy the
template to a new name. Read `references/templates.md` first. It covers the
files, the fields layouts receive, CSS scoping, and the rules that keep a
template working at any page size.

## Pitfalls

- **One heading per slide**, at the top. A heading further down stays in the
  body, and `.Title` is then empty.
- **`---` in the body** is a horizontal rule, not a slide break. Every slide
  is its own file; inkpress has no slide separator.
- **Frontmatter must be the first line.** A blank line before `---` makes it
  body text.
- **Quote YAML values** that contain `:` or start with `*`, `#`, `[` or `{`:
  `subtitle: "Q3: what changed"`.
- **Filenames sort as text**: `10-x.md` comes before `2-x.md`. Always
  zero-pad (`02-`, `10-`).
- **Unknown layout or theme names are errors.** Run `inkpress list` to see what
  each slide uses. The error message lists the valid names.
- **Notes are not output.** inkpress has no notes export. If the user wants
  notes in the PDF, ask; that needs a template change.
