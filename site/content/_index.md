---
title: inkpress
layout: hextra-home
---

{{< hextra/hero-badge >}}
  <div class="hx:w-2 hx:h-2 hx:rounded-full hx:bg-primary-400"></div>
  <span>Markdown in, PDF out</span>
{{< /hextra/hero-badge >}}

<div class="hx:mt-6 hx:mb-6">
{{< hextra/hero-headline >}}
  Slides from Markdown,&nbsp;<br class="hx:sm:block hx:hidden" />designed by a template pack
{{< /hextra/hero-headline >}}
</div>

<div class="hx:mb-12">
{{< hextra/hero-subtitle >}}
  inkpress turns a folder of Markdown slides into a PDF.&nbsp;<br class="hx:sm:block hx:hidden" />A template pack of HTML layouts and CSS themes sets the design.
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx:mb-6">
{{< hextra/hero-button text="Get started" link="docs/getting-started" >}}
</div>

{{< slides >}}
{{< slide name="title" caption="title" >}}
{{< slide name="boxes" caption="content, with boxes and a callout" >}}
{{< slide name="statement" caption="statement" >}}
{{< /slides >}}

{{< hextra/feature-grid >}}
  {{< hextra/feature-card
    title="One file per slide"
    subtitle="Each slide is a Markdown file in slides/. Slides run in filename order, and each one diffs and reviews on its own."
    link="docs/slides"
  >}}
  {{< hextra/feature-card
    title="Twelve layouts"
    subtitle="The default template has title, section, agenda, content, columns, grid, statement, image, quote and blank layouts, in light and dark."
    link="docs/layouts"
  >}}
  {{< hextra/feature-card
    title="Boxes and callouts"
    subtitle="Wrap Markdown in `:::` lines to make boxes, cards, stats and figures. GitHub's [!NOTE] syntax makes callouts."
    link="docs/boxes"
  >}}
  {{< hextra/feature-card
    title="Template packs"
    subtitle="Layouts are Go templates and themes are CSS custom properties. A deck can mix templates and themes slide by slide."
    link="docs/templates"
  >}}
  {{< hextra/feature-card
    title="Diagrams with inkline"
    subtitle="A code block of D2 with an inkline header becomes an SVG diagram that follows the slide's light or dark theme."
    link="docs/diagrams"
  >}}
  {{< hextra/feature-card
    title="Speaker notes"
    subtitle="Notes follow a `!--` line in the slide's file. They stay next to the slide and never reach the PDF."
    link="docs/slides#speaker-notes"
  >}}
  {{< hextra/feature-card
    title="Your browser prints it"
    subtitle="inkpress is a single binary. It prints with the Chrome, Chromium, Edge or Brave that is already installed."
    link="docs/cli"
  >}}
  {{< hextra/feature-card
    title="An agent skill"
    subtitle="The inkpress skill teaches Claude to plan a deck, write the slides and notes, build the PDF and change templates."
    link="docs/agent-skill"
  >}}
{{< /hextra/feature-grid >}}
