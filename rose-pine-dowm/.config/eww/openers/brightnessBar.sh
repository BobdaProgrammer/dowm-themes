#/bin/sh

if [[ "$(eww active-windows)" != *brightness-bar* ]]; then
    eww open brightness-bar
fi

bright=$(brightnessctl get)
sleep 1
if [[ "$(brightnessctl get)" == "$bright" ]];then
    eww close brightness-bar
fi
