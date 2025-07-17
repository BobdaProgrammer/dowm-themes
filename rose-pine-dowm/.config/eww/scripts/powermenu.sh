#!/bin/sh


if [[ $2 == "CONFIRMED" ]]; then 
    eww close confirm
    if [[ $1 == "--shutdown" ]]; then
        systemctl poweroff
    elif [[ $1 == "--reboot" ]]; then
        systemctl reboot
    elif [[ $1 == "--suspend" ]]; then
        playerctl --all-players pause
        amixer set Master mute
        systemctl suspend
    elif [[ $1 == "--logout" ]]; then
        kill $(pgrep -o doWM)
    elif [[ $1 == "--lock" ]]; then
        ~/lock.sh
    elif [[ $1 == "--hibernate" ]]; then
        playerctl --all-players pause
        amixer set Master mute
        systemctl hibernate
    fi
else
    eww update COMMAND="scripts/powermenu.sh $1 CONFIRMED"
    eww close logout
    eww close logout-closer
    eww open confirm
fi
