---
title: Layouts
weight: 4
---

The default template has twelve layouts and two themes, `light` and `dark`.
Each example below is a slide of the sample pack that `inkpress new` writes,
built by inkpress. Switch the site between light and dark to see both
themes.

A slide picks its layout in its frontmatter. A slide that names no layout uses
`content`.

## title

The opening slide. It shows the heading, the `subtitle`, and the author and
date from settings.toml. It has no header or footer.

{{< slide name="title" caption="title" >}}

```markdown
---
layout: title
subtitle: Markdown in, slides out
---

# Inkpress
```

## agenda

The parts of the talk as big numbered items. `current: N` highlights item N.
Repeat the slide before each part.

{{< slide name="agenda" caption="agenda, with current: 2" >}}

## section

A divider between parts, on an accent background. `kicker` is the small text
above the heading.

{{< slide name="section" caption="section" >}}

## content

The default: a heading and a body.

{{< slide name="why" caption="content" >}}

## two-column

Comparisons, or text beside code or an image. A `|||` line splits the
columns.

{{< slide name="two-column" caption="two-column" >}}

## grid

Two to eight parallel points, each in a card. `|||` lines split the cards,
and each card starts with a `###` heading.

{{< slide name="grid" caption="grid" >}}

```markdown
---
layout: grid
---

# Three ways to use a template

### Keep it

`inkpress new` writes the default template into the pack.

|||

### Change it

Edit the CSS and HTML in `.templates/default/`, or add a theme.

|||

### Swap it

Copy another template pack into `.templates/` and name it in settings.toml.
```

## statement

One big claim or number, centred. `**bold**` in the heading takes the accent
colour.

{{< slide name="statement" caption="statement" >}}

## image and image-left

Text on one half, and the frontmatter `image` filling the other.

{{< slides >}}
{{< slide name="image" caption="image" >}}
{{< slide name="image-left" caption="image-left" >}}
{{< /slides >}}

## image-full

A picture that fills the slide, with the heading and a line of text on a
scrim at the bottom. Keep the text short.

{{< slide name="image-full" caption="image-full" >}}

## quote

One large quotation, a `>` blockquote, with `cite`, `role` and an optional
`image` beside it.

{{< slide name="quote" caption="quote, with theme: dark in its frontmatter" >}}

```markdown
---
layout: quote
theme: dark
image: ../assets/cover.svg
cite: A presenter
role: after the last-minute edit
---

> I changed one word and the whole deck rebuilt in a second.
```

## blank

Full-bleed content with no padding, header or footer.

{{< slide name="blank" caption="blank, with a gradient background" >}}

## Themes

settings.toml sets the deck's theme, and a slide can set its own with
`theme: dark`, for example for a closing slide or a quote.

{{< slides >}}
{{< slide name="boxes" caption="light" theme="light" >}}
{{< slide name="boxes" caption="dark" theme="dark" >}}
{{< /slides >}}
