#!/bin/bash

status=$(playerctl status 2>/dev/null)

case "$status" in
    Playing)
        echo "󰏥"   # nf-fa-play
        ;;
    Paused)
        echo ""   # nf-md-pause_circle
        ;;
    *)
        echo ""    # Nothing or unsupported status
        ;;
esac
