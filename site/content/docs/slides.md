---
title: Slides
weight: 3
---

A slide is one Markdown file. Frontmatter at the top picks the layout and
theme, the body is Markdown, and speaker notes go at the end.

```markdown {filename="slides/05-compare.md"}
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

## Frontmatter

Frontmatter is YAML between `---` lines, or TOML between `+++` lines. It must
start on the first line of the file. inkpress reads these keys:

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

Quote YAML values that contain `:` or start with `*`, `#`, `[` or `{`:
`subtitle: "Q3: what changed"`.

An unknown layout or theme is an error. The message names the slide and lists
the valid names.

## The heading

The slide's leading `#` or `##` heading becomes `.Heading`. Layouts place it
apart from the body, so start each slide with one heading. A heading further
down stays in the body.

`---` in the body is a horizontal rule, not a slide break. Every slide is its
own file.

## Columns

Lines that hold only `|||` split the body into columns. The `two-column` and
`grid` layouts place them side by side.

{{< slide name="two-column" caption="two-column: the code block and the text are split by a ||| line" >}}

## Speaker notes

Everything after a line that holds only `!--` is the slide's notes. Notes are
not rendered. A `!--` line inside a code block is part of the code.

## Markdown

Markdown is GitHub-flavoured: tables, task lists, strikethrough and
autolinks. Raw HTML is allowed. Code blocks with a language are highlighted
in the template's Chroma style for the slide's theme.

{{< slide name="code" caption="Code and a table, in the content layout" >}}

## Images

Put images in `assets/`. A path resolves against the slide first, then the
pack, so `../assets/chart.svg` and `assets/chart.svg` both work. The first
form also works in an editor's Markdown preview.

A missing image is a warning, not an error.

## Fitting content

A slide does not scroll; text that overflows is cut off. At the default size,
a `content` slide holds about:

- a heading and up to 6 bullets of one line each, or
- a heading and about 60 words of prose, or
- a heading and a code block of about 12 lines, or a table of about 6 rows.

Split anything longer across slides, or move the detail into the notes.
