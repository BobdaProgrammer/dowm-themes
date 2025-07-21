#/bin/sh

get(){
    echo "$(pactl get-sink-volume @DEFAULT_SINK@ | grep -oP '\d+%' | head -1 | tr -d '%')"
}

if [[ $1 == "--get" ]]; then
    get
fi
