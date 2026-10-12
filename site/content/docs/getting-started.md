---
title: Getting started
weight: 1
---

## Install

inkpress is a single binary. The macOS build is signed and notarised.

```sh
brew install --cask inkcheck/tap/inkpress    # macOS and Linux
```

Or with Go:

```sh
go install github.com/inkcheck/inkpress/cmd/inkpress@latest
```

inkpress prints PDFs with a Chromium-based browser that is already installed:
Chrome, Chromium, Edge, Brave or Vivaldi. It does not bundle a browser.

## Start a pack

```sh
inkpress new talk
```

This writes a sample deck to `talk/`. It has thirteen slides that show the
layouts, boxes, code, images and a diagram, and the default template in
`talk/.templates/default/`.

## Build it

```sh
inkpress pdf talk          # -> talk/talk.pdf
inkpress html talk         # -> talk/talk.html, a preview for the browser
inkpress list talk         # slides in order, with template/layout/theme
```

[Download the sample deck as a PDF](../../sample.pdf). It is the sample pack
with two extra slides.

{{< slide name="why" caption="slides/03-why.md, in the content layout" >}}

## Write a slide

Replace the sample slides with your own. Each slide is one file in `slides/`,
and the slides run in filename order.

```markdown {filename="slides/02-problem.md"}
---
layout: two-column
---

# The problem

Reviews take four rounds.

|||

Each round waits a day.

!--
Speaker notes. Everything after the !-- line stays out of the PDF.
```

## Next

- [Slides](../slides) covers the frontmatter, columns, notes and Markdown.
- [Layouts](../layouts) shows each layout of the default template.
