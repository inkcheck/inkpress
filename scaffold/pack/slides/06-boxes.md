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

!--
Callouts use GitHub's syntax: > [!NOTE], [!TIP], [!IMPORTANT], [!WARNING]
or [!CAUTION]. Text after the marker replaces the title.
