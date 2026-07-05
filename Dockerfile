FROM alpine:latest

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
    python3

WORKDIR /kbuild/
ENTRYPOINT ["make", "-j1", "V=1"]
