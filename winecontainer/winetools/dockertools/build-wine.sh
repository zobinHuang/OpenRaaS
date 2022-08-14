#!/usr/bin/env bash

. ./build.cfg

docker build --add-host=${winhq_host}:${winhq_ip} -t dcwine .