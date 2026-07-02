#!/bin/bash

qemu-system-aarch64 \
  -M virt \
  -cpu cortex-a72 \
  -m 512M \
  -kernel ./build/kernel/boot/vmlinuz-virt \
  -initrd ./build/initrd.cpio.gz \
  -append "console=ttyAMA0 init=/init" \
  -nographic
