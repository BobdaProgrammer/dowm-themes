#!/usr/bin/env bash

## Author  : Aditya Shakya (adi1090x)
## Github  : @adi1090x
#
## Applets : Brightness (brightnessctl version)

# Import Current Theme
source "$HOME"/.config/rofi/applets/shared/theme.bash
theme="$type/$style"

# Brightness Info
# Get current brightness and max brightness from brightnessctl
backlight=$(brightnessctl get)
max_brightness=$(brightnessctl max)

# Calculate percentage brightness
backlight_percent=$(( 100 * backlight / max_brightness ))

# Get device name (assuming first backlight device)
card=$(brightnessctl --list-devices | head -n1)

if [[ $backlight_percent -ge 0 ]] && [[ $backlight_percent -le 29 ]]; then
    level="Low"
elif [[ $backlight_percent -ge 30 ]] && [[ $backlight_percent -le 49 ]]; then
    level="Optimal"
elif [[ $backlight_percent -ge 50 ]] && [[ $backlight_percent -le 69 ]]; then
    level="High"
elif [[ $backlight_percent -ge 70 ]] && [[ $backlight_percent -le 100 ]]; then
    level="Peak"
fi

# Theme Elements
prompt="${backlight_percent}%"
mesg="Device: ${card}, Level: $level"

if [[ "$theme" == *'type-1'* ]]; then
    list_col='1'
    list_row='4'
    win_width='400px'
elif [[ "$theme" == *'type-3'* ]]; then
    list_col='1'
    list_row='4'
    win_width='120px'
elif [[ "$theme" == *'type-5'* ]]; then
    list_col='1'
    list_row='4'
    win_width='425px'
elif [[ ( "$theme" == *'type-2'* ) || ( "$theme" == *'type-4'* ) ]]; then
    list_col='4'
    list_row='1'
    win_width='550px'
fi

# Options
layout=$(grep 'USE_ICON' ${theme} | cut -d'=' -f2)
if [[ "$layout" == 'NO' ]]; then
    option_1=" Increase"
    option_2=" Optimal"
    option_3=" Decrease"
    option_4=" Settings"
else
    option_1=""
    option_2=""
    option_3=""
    option_4=""
fi

# Rofi CMD
rofi_cmd() {
    rofi -theme-str "window {width: $win_width;}" \
        -theme-str "listview {columns: $list_col; lines: $list_row;}" \
        -theme-str 'textbox-prompt-colon {str: "";}' \
        -dmenu \
        -p "$prompt" \
        -mesg "$mesg" \
        -markup-rows \
        -theme ${theme}
}

# Pass variables to rofi dmenu
run_rofi() {
    echo -e "$option_1\n$option_2\n$option_3\n$option_4" | rofi_cmd
}

# Execute Command
run_cmd() {
    if [[ "$1" == '--opt1' ]]; then
        brightnessctl set +5%
    elif [[ "$1" == '--opt2' ]]; then
        brightnessctl set 25%
    elif [[ "$1" == '--opt3' ]]; then
        brightnessctl set 5%-
    elif [[ "$1" == '--opt4' ]]; then
        xfce4-power-manager-settings
    fi
}

# Actions
chosen="$(run_rofi)"
case ${chosen} in
    $option_1)
        run_cmd --opt1
        ;;
    $option_2)
        run_cmd --opt2
        ;;
    $option_3)
        run_cmd --opt3
        ;;
    $option_4)
        run_cmd --opt4
        ;;
esac
