---
title: Packs
weight: 2
---

A presentation pack is a folder:

```
talk/
  settings.toml           title, author, template, theme, page size
  README.md
  slides/                 one slide per file, in filename order: 01-intro.md, 02-…
  assets/                 images the slides use (jpg, png, svg)
  .templates/default/     a template pack; a deck can hold several
```

## settings.toml

```toml {filename="settings.toml"}
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

| Key | Meaning |
|---|---|
| `title`, `author`, `date` | shown by the templates, for example on the title slide and in the footer |
| `template` | the template pack in `.templates/`; default `default` |
| `theme` | a theme of that template; default from its `template.toml` |
| `output` | the PDF path, relative to the pack |
| `page.width` | the slide width in CSS pixels; default 1280 |
| `page.aspect`, `page.height` | the slide height, as a ratio or in pixels; default 16:9 |
| `meta` | free-form values for templates |

## Slide order

Slides run in filename order, and filenames sort as text. Zero-pad the
numbers, `02-` and `10-`, so that `10-x.md` does not come before `2-x.md`.
Leave gaps, such as `10-`, `20-`, to insert slides later.

To leave a slide out without deleting it, set `skip: true` in its
frontmatter.

## Page size

The default slide is 1280 by 720 CSS pixels. The default template sizes text
and spacing as a share of the slide, so the design scales with
`page.width`. A 4:3 deck needs only a change of `aspect`.
