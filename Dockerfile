FROM alpine:latest

ENV BUILD_TYPE=tools

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
    go

WORKDIR /apps/
COPY ./alpha/ ./alpha/
COPY ./guarddog/ ./guarddog/
COPY ./toybox/ ./toybox/

WORKDIR /build/
COPY ./isolatedbuild.sh .
COPY ./kernel.config ./.config
RUN chmod +x ./isolatedbuild.sh
ENTRYPOINT ["./isolatedbuild.sh"]
