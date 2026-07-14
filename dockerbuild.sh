#!/bin/sh

docker build --platform linux/arm64 -t sirius:latest .
docker run --rm --platform linux/arm64 --privileged -v "./build/:/mnt/" -e BUILD_TYPE=$1 sirius:latest
