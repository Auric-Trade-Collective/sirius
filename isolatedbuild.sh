#!/bin/sh
ROOT="/mnt/sirius-root"
FS="/mnt/sirius-fs"
ZBUILD="/mnt/zbuild"

source "$HOME/.cargo/env"

function env_diagnostics {
    echo "Running environment diagnostics"
    ldd $FS/bin/coreutils
}

function test {
    echo "test:x:1000:1000::/home/test:/bin/leash" > $FS/etc/passwd
    echo "/bin/" > $FS/etc/paths
}

function initramfs {
    rm -rf $ROOT/
    mkdir $ROOT/
}

function initcorefs {
    rm -rf $FS/
    mkdir $FS/

    mkdir -p $ROOT/proc/
    mkdir -p $ROOT/sys/
    mkdir -p $ROOT/dev/
    mkdir -p $ROOT/lib/
    mkdir -p $ROOT/lib/modules/
    mkdir -p $ROOT/bin/

    mkdir -p $FS/apps/
    mkdir -p $FS/etc/
    mkdir -p $FS/etc/alpha/
    mkdir -p $FS/var/log/alpha/
    mkdir -p $FS/dev/
    mkdir -p $FS/sys/
    mkdir -p $FS/proc/

    cp /build/alpha.toml $FS/etc/alpha/alpha.toml

    echo "" > /etc/passwd
    echo "" > /etc/shadow

}

function build_zfs {
    rm -rf /mnt/zbuild/
    mkdir /mnt/zbuild/

    git clone https://github.com/openzfs/zfs.git --depth=1 -b zfs-2.4.2
    cd zfs

    ./autogen.sh
    make distclean || true

    ./configure \
        --host=aarch64-linux-gnu \
        --with-linux=/mnt/kbuild/ \
        --with-linux-obj=/mnt/kbuild/ \
        --with-config=all \
        --with-spec=no \
        LIBS="-leconf -lz"

    grep -A20 "cannot create executables" /build/zfs/config.log
    find / -name "libc.a" 2>/dev/null
    find / -name "crt1.o" -o -name "crtbeginT.o" 2>/dev/null

    if [ ! -f Makefile ]; then
        echo "Configure failed! Check config.log"
        exit 1
    fi

    # Build only the modules
    make -j1 ARCH=arm64 V=1 LDFLAGS="-all-static"
    echo ZFS FINISHED: $?

    cp -a /build/zfs /mnt/zbuild/
}

function build_kernel {
    rm -rf /mnt/kbuild/
    mkdir /mnt/kbuild/

    rm -rf /mnt/kernel/
    mkdir /mnt/kernel/
    mkdir /mnt/kernel/boot/

    cd /mnt/kbuild/
    git clone --branch v7.0 https://github.com/torvalds/linux . --depth 1
    cp /build/.config ./.config
    ls -a

    make ARCH=arm64 CROSS_COMPILE=aarch64-none-elf- -j1 V=1 2>&1 | tee ./build.log
    echo "EXIT: $?"
    cp -r ./arch/arm64/boot/** /mnt/kernel/boot/

    make modules_prepare
}

function build_tools {
    cd /apps/

    echo "Building Sirius..."
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C /apps/alpha --ldflags="-s -w" -o $ROOT/init .
    chmod +x $ROOT/init

    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C /apps/guarddog --ldflags="-s -w" -o $FS/bin/guarddog .
    chmod +x $FS/bin/guarddog

    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C /apps/leash --ldflags="-s -w" -o $FS/bin/leash .
    chmod +x $FS/bin/leash

    git clone https://github.com/uutils/coreutils
    cd coreutils
    cargo build --release \
                --features "rm cat ls cp mv"
    cp ./target/release/coreutils $FS/bin/

    echo "Symlinking coreutils..."
    cd $FS/bin/
    ln -s coreutils ls
    ln -s coreutils cp
    ln -s coreutils mv
    ln -s coreutils cat
    ln -s coreutils rm
}

function build_qcow {
    cd /mnt/
    rm -f ./sirius_fs.qcow2

    rm -rf /mnt/bootstrapfs/
    mkdir /mnt/bootstrapfs/
    mkdir /mnt/bootstrapfs/lib/
    mkdir /mnt/bootstrapfs/lib/modules/
    mkdir /mnt/bootstrapfs/bin/
    mkdir /mnt/bootstrapfs/dev/
    mkdir /mnt/bootstrapfs/proc/
    mkdir /mnt/bootstrapfs/sys/
    cp /mnt/zbuild/zfs/module/zfs.ko /mnt/bootstrapfs/lib/modules/
    cp /mnt/zbuild/zfs/module/spl.ko /mnt/bootstrapfs/lib/modules/
    cp /mnt/zbuild/zfs/zpool /mnt/bootstrapfs/bin/zpool
    cp /mnt/zbuild/zfs/zfs /mnt/bootstrapfs/bin/zfs
    cp -a $FS/ /mnt/bootstrapfs/

    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C /apps/zfsbootstrap --ldflags="-s -w" -o /mnt/bootstrapfs/bin/init .
    chmod +x /mnt/bootstrapfs/bin/init

    cd /mnt/bootstrapfs/
    find . -print0 | cpio --null -ov --format=newc > ../initrdboot.cpio
    cd ../

    rm initrdboot.cpio.gz
    gzip initrdboot.cpio

    qemu-img create -f qcow2 /mnt/sirius_fs.qcow2 2G
    qemu-system-aarch64 \
        -M virt \
        -cpu cortex-a72 \
        -m 512M \
        -kernel /mnt/kernel/boot/Image \
        -initrd ./initrdboot.cpio.gz \
        -drive file=./sirius_fs.qcow2,format=qcow2,if=virtio \
        -nographic \
        -no-reboot \
        -append "console=ttyAMA0,115200 earlycon=pl011,0x09000000 rdinit=/bin/init panic=-1"
}

function finalize_package {
    cp /mnt/zbuild/zfs/module/zfs.ko $ROOT/lib/modules/
    cp /mnt/zbuild/zfs/module/spl.ko $ROOT/lib/modules/
    cp /mnt/zbuild/zfs/zpool $ROOT/bin/zpool
    cp /mnt/zbuild/zfs/zfs $ROOT/bin/zfs

    cd $ROOT/
    find . -print0 | cpio --null -ov --format=newc > ../initrd.cpio
    cd ../

    rm initrd.cpio.gz
    gzip initrd.cpio

    build_qcow
}

if [ $BUILD_TYPE == "full" ]; then
    build_kernel
    build_zfs
fi


if [ $BUILD_TYPE == "zfs" ]; then
    build_zfs
fi

initramfs
initcorefs
build_tools
test
finalize_package
env_diagnostics
