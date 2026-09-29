# Diagrams

```inkline
vars: { inkline: { diagram: flowchart; design_system: carbon } }

md: Markdown slides {class: start}
inkpress: inkpress {class: process}
pdf: PDF {class: end}
md -> inkpress: render
inkpress -> pdf: print
```

???
inkline renders this block to SVG. Without inkline installed, the slide shows the D2 source.
