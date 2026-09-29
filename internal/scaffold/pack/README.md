# Inkpress deck

Build the PDF:

```sh
inkpress pdf            # -> <directory name>.pdf
inkpress html           # -> an HTML preview to open in a browser
```

## Layout

```
settings.toml          title, author, template, theme, page size
slides/                one Markdown file per slide, in filename order
assets/                images the slides use (jpg, png, svg)
.templates/default/    the template pack: layouts, themes, header, footer, fonts
```

## A slide

```markdown
---
layout: two-column      # title, section, content, two-column, image, quote, blank
theme: dark             # light or dark
---

# Heading

Left column

|||

Right column

???
Speaker notes. Everything after ??? stays out of the slides.
```

Frontmatter keys: `template`, `layout`, `theme`, `title`, `class` (extra CSS
classes), `background` (a colour, gradient or image path), `image` (for the
image layout), `header` and `footer` (`false` hides them), `skip` (`true`
leaves the slide out). Layouts can read any other key as `.Meta.<key>`.

Images resolve against the slide (`../assets/chart.svg`) or the pack
(`assets/chart.svg`).
