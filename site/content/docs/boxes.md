---
title: Boxes and callouts
weight: 5
---

{{< slide name="boxes" caption="slides/06-boxes.md: a callout, two stats and a warning box in a row" >}}

````markdown {filename="slides/06-boxes.md"}
# Boxes and callouts

> [!TIP]
> Wrap Markdown in `:::` lines to put it in a box. The words after `:::` are
> CSS classes that the template styles.

::: row
::: stat
# 1
Markdown file per slide
:::
::: stat tip
# 0
design decisions in the slide
:::
::: box warning
### Watch out
A slide does not scroll. Text that overflows is cut off.
:::
:::
````

## Boxes

Markdown between `:::` lines goes in a div. `::: box note` opens
`<div class="box note">`, and a bare `:::` closes it. Boxes nest, and the
template styles the classes.

The default template has these classes:

| Class | Use |
|---|---|
| `box` | a shaded panel with an accent edge; a leading `###` becomes its label |
| `card` | a shaded panel with an accent top; the `grid` layout wraps each column in one |
| `row` | puts the blocks inside it side by side, equal widths |
| `stat` | a big number (a heading) over a short line of text |
| `figure` | an image, then a caption paragraph |
| `gallery` | images in a row, cropped to the same shape |
| `note`, `tip`, `important`, `warning`, `caution` | the tone of a `box`, `card` or `stat` |
| `lead`, `muted`, `center` | larger intro text, grey text, centred text |

A box cannot span a `|||` line. One left open closes at the end of its
column, so put boxes inside columns instead.

Prefer a layout, such as `grid` or `statement`, over rebuilding one with
boxes.

## Callouts

Callouts use GitHub's syntax: a blockquote whose first line is `[!NOTE]`,
`[!TIP]`, `[!IMPORTANT]`, `[!WARNING]` or `[!CAUTION]`.

```markdown
> [!WARNING]
> The PDF needs a Chromium-based browser.
```

This becomes `<blockquote class="callout warning">`, whose first child is a
`.callout-title` paragraph. Text after the marker, as in `> [!TIP] Try this`,
replaces the title.

Use at most one callout per slide.
