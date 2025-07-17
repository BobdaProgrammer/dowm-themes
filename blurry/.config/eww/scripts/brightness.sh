#!/bin/sh

val="$2"

get(){
    echo "$(brightnessctl get)"
}

get_percent(){ 
    per="$(/home/sam/.config/eww/scripts/percent.sh $(brightnessctl get) 255)"
    echo "$(awk "BEGIN { print int($per + 0.5) }")"
}

setBrightness(){
    brightnessctl set "$val"
}

decrease(){
    brightnessctl set 5-
}

increase(){
    brightnessctl set +5
}

if [[ $1 == "--get" ]]; then
    get
elif [[ $1 == "--set" ]]; then
    setBrightness
elif [[ $1 == "--get-percent" ]]; then
    get_percent
elif [[ $1 == "--increase" ]]; then
    increase
elif [[ $1 == "--decrease" ]]; then
    decrease
fi
