#!/bin/sh

docker build -t sirius:latest .
docker run --rm -v "./build/:/mnt/" -e BUILD_TYPE=$1 sirius:latest
