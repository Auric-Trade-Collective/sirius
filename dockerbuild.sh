#!/bin/bash

function initcorefs {
    mkdir ./build/sirius-root/

    ROOT="./build/sirius-root"

    mkdir -p $ROOT/etc/alpha/
    mkdir -p $ROOT/proc/
    mkdir -p $ROOT/sys/
    mkdir -p $ROOT/dev/
    mkdir -p $ROOT/var/log/alpha/

    mkdir -p $ROOT/apps/

    mknod -m 600 $ROOT/dev/console c 5 1
    mknod -m 666 $ROOT/dev/null c 1 3
    mknod -m 666 $ROOT/dev/zero c 1 5

cat <<EOF > $ROOT/etc/alpha/alpha.toml
[isolated]
# Rules for containers (namespaces, cgroups, etc)

[host]
# Rules for core services (compositor, init-bridge)
ls = {
    name = "/bin/ls",
    args = "/"
}
EOF

    # Inject packages

    # mkdir -p $ROOT/apps/fun/
    # mkdir -p $ROOT/apps/fun/data/
    # GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C ./fun --ldflags="-s -w" -o ../build/sirius-root/apps/fun/fun .
    # cp ./fun/package.toml $ROOT/apps/fun/
}

function buildinit {
    mkdir ./build/sirius-root/

    echo "Building Sirius..."
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C ./alpha --ldflags="-s -w" -o ../build/sirius-root/init .
    chmod +x ./build/sirius-root/init

    initcorefs

    cd ./build/sirius-root/
    find . -print0 | cpio --null -ov --format=newc > ../initrd.cpio
    cd ../

    gzip initrd.cpio

    cd $root
}

function full_reset {
    rm -rf ./build/
    mkdir ./build/sirius-root/
    mkdir ./build/kernelb/
    mkdir ./build/kernel/boot/
}

function reset_tools {
    rm -rf ./build/sirius-root/
    rm ./build/initrd.cpio.gz
    mkdir ./build/sirius-root/
}

root=$(pwd)

if [ ! -z "$1"  ] && [ "$1" == "full" ]; then
    full_reset

    git clone https://github.com/torvalds/linux ./build/kernelb/ --depth 1
    cp kernel.config "./build/kernelb/.config"
    docker build -t sirius .
    docker run --rm -v "./build/kernelb/:/kbuild/" sirius:latest

    mkdir ./build/kernel/boot/
    cp ./build/kernelb/arch/arm64/boot/** ./build/kernel/boot/

    buildinit
    echo "Done!"
fi


if [ ! -z "$1"  ] && [ "$1" == "tools" ]; then
    reset_tools
    buildinit

    echo "Done!"
fi
