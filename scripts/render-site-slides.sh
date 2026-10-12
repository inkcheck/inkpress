#!/bin/sh
# Builds the sample pack, plus the slides in site/slides/, in both themes, and
# writes each page as a PNG into the site's static files. Needs a
# Chromium-based browser and pdftoppm (poppler); inkline renders the diagram
# slide when it is on the PATH.
set -eu
cd "$(dirname "$0")/.."
go build -o bin/inkpress ./cmd/inkpress
out=site/static/slides
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

rm -rf "$out"
mkdir -p "$out"
./bin/inkpress new "$work/sample" >/dev/null
cp site/slides/*.md "$work/sample/slides/"

for theme in light dark; do
  pack="$work/sample-$theme"
  cp -R "$work/sample" "$pack"
  sed "s/^theme = \"light\"/theme = \"$theme\"/" "$work/sample/settings.toml" > "$pack/settings.toml"
  ./bin/inkpress pdf "$pack" -o "$work/$theme.pdf" >/dev/null
  pdftoppm -png -r 96 "$work/$theme.pdf" "$work/$theme"
  # Name each page after its slide file: 03-why.md -> why-light.png. The
  # list leaves out skipped slides, as the PDF does. pdftoppm zero-pads the
  # page numbers, so the first page left is the next.
  for f in $(./bin/inkpress list "$pack" | awk '{ print $2 }'); do
    page=$(ls "$work/$theme"-*.png | head -n 1)
    name=$(basename "$f" .md | sed 's/^[0-9]*-//')
    mv "$page" "$out/$name-$theme.png"
  done
done
cp "$work/light.pdf" site/static/sample.pdf
echo "rendered $(ls "$out" | wc -l | tr -d ' ') slides into $out"
