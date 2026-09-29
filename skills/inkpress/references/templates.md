# inkpress template packs

A template pack is a directory in the deck's `.templates/`. A slide uses the
template named in its frontmatter, or in `settings.toml`.

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

The quickest start is to copy `.templates/default` to a new name and edit it.

## template.toml

```toml
description = "Clean, type-led slides"
layout = "content"     # used by slides that name no layout; must exist
theme = "light"        # used when neither slide nor settings name one

[code_styles]          # Chroma style per theme; default "github"
light = "github"
dark = "github-dark"
```

Chroma style names: https://xyproto.github.io/splash/docs/. An unknown name
is an error.

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

```html
{{.Header}}
<div class="frame">
  {{with .Heading}}<div class="heading">{{.}}</div>{{end}}
  <div class="columns">
    {{range .Columns}}<div class="column">{{.}}</div>{{end}}
  </div>
</div>
{{.Footer}}
```

| Field | |
|---|---|
| `.Heading` | the slide's leading `#`/`##` heading, as HTML; empty if none |
| `.Body` | everything after the heading, as HTML |
| `.Content` | `.Heading` and `.Body` together |
| `.Columns` | `.Body` split at `\|\|\|` lines; one item when there are none |
| `.Title` | frontmatter `title`, or the heading's text |
| `.Number`, `.Total` | this slide's number from 1, and the slide count |
| `.Template`, `.Layout`, `.Theme` | the names in use |
| `.Image` | frontmatter `image`, resolved to a URL |
| `.Meta` | the whole frontmatter: `{{.Meta.subtitle}}` |
| `.Header`, `.Footer` | rendered partials; empty if the slide sets `header: false` / `footer: false` |
| `.Deck.Title`, `.Deck.Author`, `.Deck.Date`, `.Deck.Meta` | from settings.toml |

Functions: `asset "logo.svg"` returns the URL of `assets/logo.svg` in the
template. Use it in `src` and `href` attributes. The standard template
functions (`with`, `if`, `range`, `and`, `or`, `eq`, `printf`) work.

The layout decides whether the header and footer appear. The default `title`
layout leaves out both. Partials get the same fields, with `.Header` and
`.Footer` empty.

## CSS

inkpress wraps `style.css` in `@scope (.tpl-<name>)` and each theme in
`@scope (.tpl-<name>.theme-<theme>)`. The scope root is `.page`, so write
selectors for `.slide` and its descendants:

```css
.slide { font-family: "Inter", sans-serif; color: var(--fg); background: var(--bg); }
.layout-title h1 { font-size: 9cqh; }
```

- **Size in `cqh` and `cqw`**: 1% of the slide's height or width. The design
  then scales when settings.toml changes the page size. Avoid `px` for text
  and spacing. Avoid `vh`/`vw`, which refer to the browser window.
- **Themes set custom properties only**, on `.slide`: `--bg`, `--fg`,
  `--accent` and so on. `style.css` uses the properties. Adding a theme is
  then one small file.
- **Theme files come after `style.css`** and have the same specificity for
  `.slide`. A layout that overrides a theme property must be more specific:
  `.slide.layout-section { --fg: white; }`.
- **`@font-face` and `@import`** may sit anywhere in `style.css`; inkpress moves
  them to the top level. Relative `url()`s resolve against the CSS file, so
  `url(fonts/Inter.woff2)` and `url(assets/bg.svg)` work.
- **Colour a single-colour logo with the theme** with a mask:
  `.logo { background: var(--accent); mask: url(assets/logo.svg) center / contain no-repeat; }`.
- **Code blocks** are `<pre class="chroma">` with Chroma classes. The
  template's `pre` rule can set a background with `!important` to override
  the Chroma style's background.
- **Blocks** from `::: a b` lines are `<div class="a b">`, and callouts are
  `<blockquote class="callout <kind>">` whose first child is
  `<p class="callout-title">`. The kind is the lower-cased marker (`note`,
  `tip`, …) and any other word works too. Style every class the skill
  documents (`box`, `card`, `row`, `stat`, `figure`, `gallery`, the five
  tones, `lead`, `muted`, `center`) so decks keep working when they switch
  templates. The default template gives each tone a theme property
  (`--note`, `--tip`, …) and sets `--tone` from it.
- **Inkline diagrams** are `<figure class="diagram"><img …></figure>`. Give the
  figure `flex: 1; min-height: 0` in a column flexbox, and the img
  `max-height: 100%; object-fit: contain`, so diagrams fit the space left.
- **Fonts must be files in the template**, or installed on the machine.
  Fonts from web URLs may not load in time for printing. Keep the font's
  licence next to it.

## Checking a template

1. `inkpress list` shows every slide's template, layout and theme. An unknown
   name fails with the list of valid ones.
2. Build a test deck that uses every layout in every theme, with the longest
   content each layout should hold.
3. Print it with `inkpress pdf` and look at the pages. Check overflow, contrast
   in the dark theme, header and footer placement, and code blocks.
4. Change `[page] aspect` to `4:3` and print again. Nothing should overflow or
   overlap.
