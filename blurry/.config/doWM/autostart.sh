#!/bin/bash


killall dunst
# Start Eww daemon (ensure it doesn't already run)
~/.config/eww/daemon &


notify-send "startup"
feh --bg-scale ~/wallpapers/picker/trees.jpg

polybar &

#~/PICOM --experimental-backends &
~/clones/picom/build/src/picom &

/usr/lib/polkit-gnome/polkit-gnome-authentication-agent-1 &
kdeconnectd &

kdeconnect-cli --refresh -n "Pixel 3a"
