#!/bin/bash


killall dunst
# Start Eww daemon (ensure it doesn't already run)
~/.config/eww/daemon &

notify-send "startup"
feh --bg-scale ~/wallpapers/picker/brightblue.jpg

eww open dock

polybar &

#~/PICOM --experimental-backends &
~/clones/picom/build/src/picom &

