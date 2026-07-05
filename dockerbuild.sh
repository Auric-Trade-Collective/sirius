#!/bin/bash

root=$(pwd)

if [ ! -z "$1"  ] && [ "$1" == "full" ]; then
    rm -rf ./build/
    mkdir ./build/sirius-root/
    rm -rf ./build/kernelb/
    rm -rf ./build/kernel/
    mkdir ./build/kernelb/
    mkdir ./build/kernel/boot/
    git clone https://github.com/torvalds/linux ./build/kernelb/ --depth 1
    cp kernel.config "./build/kernelb/.config"
    docker build -t sirius .
    docker run --rm -v "./build/kernelb/:/kbuild/" sirius:latest

    mkdir ./build/kernel/boot/
    cp ./build/kernelb/arch/arm64/boot/** ./build/kernel/boot/

    echo "Building Sirius..."
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C ./alpha --ldflags="-s -w" -o ../build/sirius-root/init .
    chmod +x ./build/sirius-root/init

    ./initcorefs.sh

    cd ./build/sirius-root/
    find . -print0 | cpio --null -ov --format=newc > ../initrd.cpio
    cd ../

    gzip initrd.cpio

    cd $root
    echo "Done!"
fi


if [ ! -z "$1"  ] && [ "$1" == "tools" ]; then
    rm -rf ./build/sirius-root/
    rm ./build/initrd.cpio.gz

    mkdir ./build/sirius-root/

    echo "Building Sirius..."
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C ./alpha --ldflags="-s -w" -o ../build/sirius-root/init .
    chmod +x ./build/sirius-root/init

    ./initcorefs.sh

    cd ./build/sirius-root/
    find . -print0 | cpio --null -ov --format=newc > ../initrd.cpio
    cd ../

    gzip initrd.cpio

    cd $root
    echo "Done!"
fi
