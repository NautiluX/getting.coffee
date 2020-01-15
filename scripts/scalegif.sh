#!/bin/bash

# thanks to Cassidy K (https://github.com/cassidycodes) for the idea.
# original script at https://cassidy.codes/blog/2017/04/25/ffmpeg-frames-to-gif-optimization/

palette="/tmp/palette.png"

filters="fps=5,scale=320:-1:flags=lanczos"

ffmpeg -v warning -i "$1" -vf "$filters,palettegen=stats_mode=diff" -y $palette

ffmpeg -t 5 -i "$1" -i $palette -lavfi "$filters,paletteuse=dither=bayer:bayer_scale=5:diff_mode=rectangle" -y "$2"
