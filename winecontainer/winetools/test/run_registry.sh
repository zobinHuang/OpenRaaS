#!/usr/bin/env bash

docker pull registry

docker run -d -p 5000:5000 --name=registry --restart=always --privileged=true -v /usr/local/docker_registry:/var/lib/registry registry:latest