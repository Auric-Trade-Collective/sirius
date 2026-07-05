#!/bin/sh

function initcorefs {
    mkdir /mnt/sirius-root/
    ROOT="/mnt/sirius-root"

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
}

function build_kernel {
    rm -rf /kbuild/
    mkdir /kbuild/

    rm -rf /mnt/kernel/
    mkdir /mnt/kernel/
    mkdir /mnt/kernel/boot/

    cd /kbuild/
    git clone https://github.com/torvalds/linux . --depth 1
    cp /build/.config ./.config
    ls -a

    make -j1 V=1
    cp -r ./arch/arm64/boot/** /mnt/kernel/boot/
}

function build_tools {
    rm -rf /mnt/sirius-root/
    mkdir /mnt/sirius-root/

    cd /apps/

    echo "Building Sirius..."
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C /apps/alpha --ldflags="-s -w" -o /mnt/sirius-root/init .
    chmod +x /mnt/sirius-root/init
}

function finalize_package {
    cd /mnt/sirius-root/
    find . -print0 | cpio --null -ov --format=newc > ../initrd.cpio
    cd ../

    rm initrd.cpio.gz
    gzip initrd.cpio
}


if [ $BUILD_TYPE == "full" ]; then
    build_kernel
fi

build_tools
initcorefs
finalize_package
