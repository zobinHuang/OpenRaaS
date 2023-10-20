#!/usr/bin/env bash

ffplay -protocol_whitelist file,rtp,udp -i h264.sdp
