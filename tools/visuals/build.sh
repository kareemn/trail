#!/usr/bin/env bash
# Regenerate visuals/*.gif and README.md from TRAIL_visuals.html.
#
#   tools/visuals/build.sh                 # every animation, then README
#   tools/visuals/build.sh records,tamper  # only these data-scene names, then README
#
# Needs: node, python3, ffmpeg, Google Chrome. Run from anywhere.
# Tunables (defaults match the committed GIFs):
#   SLOWDOWN=1.8  how much slower than the page
#   FPS=10        GIF frame rate (10 gives exact 1/100 s GIF delays)
#   HOLD=3        seconds the last frame is held before looping
#   WIDTH=900     GIF width in pixels
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TOOLS="$ROOT/tools/visuals"
SLOWDOWN="${SLOWDOWN:-1.8}" FPS="${FPS:-10}" HOLD="${HOLD:-3}" WIDTH="${WIDTH:-900}"
ONLY="${1:-}"

[ -d "$TOOLS/node_modules" ] || (cd "$TOOLS" && npm ci --silent)

FRAMES="$(mktemp -d)"
trap 'rm -rf "$FRAMES"' EXIT

node "$TOOLS/capture.js" "$ROOT/TRAIL_visuals.html" "$FRAMES" "$FPS" "$SLOWDOWN" $ONLY

mkdir -p "$ROOT/visuals"
# A full run replaces every GIF, so renamed or removed animations leave nothing behind.
[ -z "$ONLY" ] && rm -f "$ROOT"/visuals/*.gif

for dir in "$FRAMES"/*/; do
  name="$(basename "$dir")"
  ffmpeg -loglevel error -y -framerate "$FPS" -i "$dir/%04d.png" \
    -filter_complex "scale=${WIDTH}:-1:flags=lanczos,tpad=stop_mode=clone:stop_duration=${HOLD},split[a][b];[a]palettegen=max_colors=128:stats_mode=diff[p];[b][p]paletteuse=dither=none:diff_mode=rectangle" \
    -loop 0 "$ROOT/visuals/$name.gif"
  echo "visuals/$name.gif  $(du -h "$ROOT/visuals/$name.gif" | cut -f1)"
done

python3 -I "$TOOLS/page2md.py" "$ROOT/TRAIL_visuals.html" visuals "$ROOT/README.md"
