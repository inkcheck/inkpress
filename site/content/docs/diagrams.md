---
title: Diagrams
weight: 6
---

A ` ```inkline ` code block, or a ` ```d2 ` block with an inkline header, is
rendered to SVG by [inkline](https://inkcheck.github.io/inkline/). The
diagram is scaled to fit the space the slide has left.

{{< slide name="diagram" caption="slides/11-diagram.md: an inkline flowchart in the Carbon design system" >}}

````markdown {filename="slides/11-diagram.md"}
# Diagrams

```inkline
vars: { inkline: { diagram: flowchart; design_system: carbon } }

md: Markdown slides {class: start}
inkpress: inkpress {class: process}
pdf: PDF {class: end}
md -> inkpress: render
inkpress -> pdf: print
```
````

## Themes

Light slides get light diagrams and dark slides get dark ones. Leave `theme`
out of the inkline header so that the diagram follows the slide. A header
that sets a theme keeps it.

## Installing inkline

inkpress runs the `inkline` on the PATH. To use another executable, set
`INKPRESS_INKLINE` to its path.

```sh
brew install --cask inkcheck/tap/inkline
```

Without inkline, the slide shows the D2 source and inkpress warns you. The
build does not stop.

## Fitting diagrams

Prefer inkline's default top-to-bottom direction. A wide diagram shrinks to
fit the slide's width and becomes hard to read.

` ```verf ` blocks and verf headers from older decks still work.
