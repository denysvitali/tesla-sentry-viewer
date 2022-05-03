#!/bin/bash

MAIN_DIR="output-2"
#   -hwaccel vaapi -hwaccel_output_format vaapi \
#   -vaapi_device /dev/dri/renderD128 \
#   -codec:v h264_vaapi \

ffmpeg \
  -i "$MAIN_DIR/left_repeater.mp4" \
  -i "$MAIN_DIR/front.mp4" \
  -i "$MAIN_DIR/right_repeater.mp4" \
  -i "$MAIN_DIR/back.mp4" \
  -filter_complex \
  "[0:v][1:v][2:v][3:v]xstack=inputs=4:layout=0_0|w0_0|w0+w1_0|w0_h0[out]" \
  -map "[out]" -ac 2 \
  "$MAIN_DIR/merged.mp4"
