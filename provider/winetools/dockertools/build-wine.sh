#!/usr/bin/env bash

. ./build.cfg

# normal
# docker build --add-host=${winhq_host}:${winhq_ip} -t dcwine -f Dockerfile .
# with cuda
docker build --add-host=${winhq_host}:${winhq_ip} -t dcwine_nvidia -f Dockerfile_nvidia .