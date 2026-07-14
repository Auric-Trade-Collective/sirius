#!/bin/bash

# qemu-system-aarch64 \
#   -M virt \
#   -cpu cortex-a72 \
#   -m 512M \
#   -kernel ./build/kernel/boot/vmlinux \
#   -initrd ./build/initrd.cpio.gz \
#   -append "console=ttyAMA0 init=/init" \
#   -nographic

# qemu-system-aarch64 \
#   -M virt \
#   -cpu cortex-a72 \
#   -m 512M \
#   -kernel ./build/kernel/boot/Image \
#   -initrd ./build/initrd.cpio.gz \
#   -append "console=ttyAMA0,115200 earlycon=pl011,0x09000000 rdinit=/init loglevel=8 ignore_loglevel" \
#   -nographic

qemu-system-aarch64 \
  -M virt \
  -cpu cortex-a72 \
  -m 512M \
  -kernel ./build/kernel/boot/Image \
  -initrd ./build/initrd.cpio.gz \
  -drive file=./build/sirius_fs.qcow2,format=qcow2,if=virtio \
  -append "console=ttyAMA0,115200 earlycon=pl011,0x09000000 rdinit=/init loglevel=8 ignore_loglevel" \
  -nographic
