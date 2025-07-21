
#!/bin/bash

DIR="$HOME/wallpapers/Downloads"

for img in "$DIR"/*.{jpg,jpeg,png,gif}; do
  [ -f "$img" ] || continue
  # Escape any single quotes or spaces in path for safe embedding
  escaped_img=$(printf '%q' "$img")
  echo "  (box :class \"wallpaper-thumb\" :style \"background-image: url('$escaped_img');\")"
done
