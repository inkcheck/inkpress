---
title: Templates
weight: 7
---

A template pack is a directory in the deck's `.templates/`. A slide uses the
template named in its frontmatter, or in settings.toml. The quickest start is
to copy `.templates/default` to a new name and edit it.

```
.templates/<name>/
  template.toml          description, default layout and theme, code styles
  style.css              styles for .slide and everything inside it
  themes/<theme>.css     custom properties on .slide, one file per theme
  layouts/<layout>.html  Go html/template, one per layout
  partials/header.html   optional; rendered into .Header
  partials/footer.html   optional; rendered into .Footer
  partials/<other>.html  optional; call with {{template "partial/<other>" .}}
  assets/                logos, backgrounds
  fonts/                 font files
```

To add the default template to an existing pack, run
`inkpress new --template-only`.

## template.toml

```toml {filename="template.toml"}
description = "Clean, type-led slides"
layout = "content"     # used by slides that name no layout; must exist
theme = "light"        # used when neither slide nor settings name one

[code_styles]          # Chroma style per theme; default "github"
light = "github"
dark = "github-dark"
```

Chroma style names are listed at
[xyproto.github.io/splash](https://xyproto.github.io/splash/docs/). An
unknown name is an error.

## What inkpress generates

```html
<div class="page tpl-<template> theme-<theme>">
  <section class="slide layout-<layout> <frontmatter class>" style="background: …">
    … the layout's output …
  </section>
</div>
```

Each `.page` is exactly the slide size, a CSS size container, and one PDF
page. `.slide` fills it, and `overflow: hidden` clips anything that does not
fit.

## Layouts

A layout is an `html/template` file. It produces the inside of `.slide`.

```html {filename="layouts/two-column.html"}
{{.Header}}
<div class="frame">
  {{with .Heading}}<div class="heading">{{.}}</div>{{end}}
  <div class="columns">
    {{range .Columns}}<div class="column">{{.}}</div>{{end}}
  </div>
</div>
{{.Footer}}
```

Layouts and partials get these fields:

| Field | |
|---|---|
| `.Heading` | the slide's leading `#` or `##` heading, as HTML; empty if none |
| `.Body` | everything after the heading, as HTML |
| `.Content` | `.Heading` and `.Body` together |
| `.Columns` | `.Body` split at `\|\|\|` lines; one item when there are none |
| `.Title` | frontmatter `title`, or the heading's text |
| `.Number`, `.Total` | this slide's number from 1, and the slide count |
| `.Template`, `.Layout`, `.Theme` | the names in use |
| `.Image` | frontmatter `image`, resolved to a URL |
| `.Meta` | the whole frontmatter: `{{.Meta.subtitle}}` |
| `.Header`, `.Footer` | rendered partials; empty if the slide sets `header: false` or `footer: false` |
| `.Deck.Title`, `.Deck.Author`, `.Deck.Date`, `.Deck.Meta` | from settings.toml |

The function `asset "logo.svg"` returns the URL of `assets/logo.svg` in the
template. Use it in `src` and `href` attributes. The standard template
functions, such as `with`, `if`, `range`, `eq` and `printf`, work.

The layout decides whether the header and footer appear. The default `title`
layout leaves out both.

## CSS

A deck can mix templates and themes, so inkpress wraps `style.css` in
`@scope (.tpl-<name>)` and each theme in `@scope (.tpl-<name>.theme-<theme>)`.
Write selectors for `.slide` and its descendants:

```css
.slide { font-family: "Inter", sans-serif; color: var(--fg); background: var(--bg); }
.layout-title h1 { font-size: 9cqh; }
```

- **Size in `cqh` and `cqw`**, 1% of the slide's height or width. The design
  then scales when settings.toml changes the page size. Avoid `px` for text
  and spacing, and avoid `vh` and `vw`, which refer to the browser window.
- **Themes set custom properties only**, on `.slide`: `--bg`, `--fg`,
  `--accent` and so on. `style.css` uses the properties, so a new theme is one
  small file.
- **Theme files come after `style.css`** with the same specificity for
  `.slide`. A layout that overrides a theme property must be more specific:
  `.slide.layout-section { --fg: white; }`.
- **`@font-face` and `@import`** may sit anywhere in `style.css`; inkpress
  moves them to the top level. Relative `url()`s resolve against the CSS file.
- **Code blocks** are `<pre class="chroma">` with Chroma classes.
- **Boxes** are `<div>`s with the classes from their `:::` line, and callouts
  are `<blockquote class="callout <kind>">`. Style every class the default
  template has, so that decks keep working when they switch templates.
- **Diagrams** are `<figure class="diagram"><img …></figure>`. Give the
  figure `flex: 1; min-height: 0` in a column flexbox, and the image
  `max-height: 100%; object-fit: contain`.
- **Fonts must be files in the template**, or installed on the machine. Fonts
  from web URLs may not load in time for printing.

```css {filename="themes/light.css"}
/* Colours for the light theme. style.css uses only these properties. */
.slide {
  --bg: #ffffff;
  --fg: #161616;
  --muted: #6f6f6f;
  --accent: #0f62fe;
  --surface: #f4f4f4;
  --border: #e0e0e0;
  --section-bg: #0f62fe;
  --section-fg: #ffffff;
  /* Callout and box tones: note, tip, important, warning, caution. */
  --note: #0f62fe;
  --tip: #24a148;
  --important: #8a3ffc;
  --warning: #b28600;
  --caution: #da1e28;
}
```

## Checking a template

1. `inkpress list` shows every slide's template, layout and theme.
2. Build a test deck that uses every layout in every theme, with the longest
   content each layout should hold.
3. Print it with `inkpress pdf` and look at the pages. Check overflow,
   contrast in the dark theme, the header and footer, and code blocks.
4. Change `[page] aspect` to `4:3` and print again. Nothing should overflow
   or overlap.
