FROM alpine:latest

ENV BUILD_TYPE=tools

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk update
RUN apk add --no-cache \
    build-base \
    bc \
    flex \
    bison \
    elfutils-dev \
    openssl-dev \
    ncurses-dev \
    git \
    rsync \
    gcc \
    make \
    perl \
    python3 \
    rust \
    go \
    qemu-img \
    qemu-aarch64 \
    qemu-arm \
    e2fsprogs \
    gcc-aarch64-none-elf \
    alpine-sdk \
    linux-headers \
    autoconf \
    automake \
    util-linux-dev \
    libtirpc-dev \
    libtool \
    bash \
    util-linux-static \
    zlib-static \
    libtirpc-static \
    musl-libintl \
    libeconf-static \
    openssl-libs-static \
    elfutils-dev \
    qemu-system-aarch64

WORKDIR /apps/
COPY ./alpha/ ./alpha/
COPY ./guarddog/ ./guarddog/
COPY ./toybox/ ./toybox/
COPY ./zfsbootstrap/ ./zfsbootstrap/

WORKDIR /build/
COPY ./isolatedbuild.sh .
COPY ./kernel.config ./.config
COPY ./prebuild/alpha.toml ./alpha.toml
RUN chmod +x ./isolatedbuild.sh
ENTRYPOINT ["./isolatedbuild.sh"]
