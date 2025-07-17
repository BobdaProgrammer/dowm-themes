
#/bin/sh

if [[ "$(eww active-windows)" != *volume-bar* ]]; then
    eww open volume-bar
fi

volume=$(pactl get-sink-volume @DEFAULT_SINK@ | grep -oP '\d+%' | head -1 | tr -d '%')
sleep 1
if [[ "$(pactl get-sink-volume @DEFAULT_SINK@ | grep -oP '\d+%' | head -1 | tr -d '%')" == "$volume" ]]; then
    eww close volume-bar
fi
