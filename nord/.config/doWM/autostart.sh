#!/bin/bash


dunst &
# Start Eww daemon (ensure it doesn't already run)

feh --bg-scale ~/wallpapers/picker/nordlookoutmountains.png

polybar &

~/clones/picom/build/src/picom &
