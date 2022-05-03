#!/bin/bash

MAIN_DIR="output"
#   -hwaccel vaapi -hwaccel_output_format vaapi \
#   -vaapi_device /dev/dri/renderD128 \
#   -codec:v h264_vaapi \

ffmpeg \
  -hwaccel vaapi -hwaccel_output_format vaapi \
  -vaapi_device /dev/dri/renderD128 \
  -i "$MAIN_DIR/$1" \
  -vf scale_vaapi=w=800:h=-1,hwdownload,format=nv12 \
  -c:v libx264 \
  -crf 35 \
  "$MAIN_DIR-2/$1"
