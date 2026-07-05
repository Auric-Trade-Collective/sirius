#!/bin/bash
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

mkdir -p $ROOT/apps/fun/
mkdir -p $ROOT/apps/fun/data/
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C ./fun --ldflags="-s -w" -o ../build/sirius-root/apps/fun/fun .
cp ./fun/package.toml $ROOT/apps/fun/
