rm -rf ./build/

base=$(pwd)

echo "Pulling alpine base..."
curl -O https://dl-cdn.alpinelinux.org/alpine/v3.24/releases/aarch64/alpine-minirootfs-3.24.1-aarch64.tar.gz

echo "Pulling kernel..."
mkdir -p ./build/kernel/boot/
curl -o ./build/kernel/boot/vmlinuz-virt \
  https://dl-cdn.alpinelinux.org/alpine/v3.24/releases/aarch64/netboot/vmlinuz-virt

echo "Extracting filesystem..."
mkdir -p ./build/sirius-root/
tar -xzf alpine-minirootfs-*.tar.gz -C ./build/sirius-root/

rm alpine-minirootfs-*.tar.gz

echo "Building Sirius..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C ./alpha --ldflags="-s -w" -o ../build/sirius-root/init .
chmod +x ./build/sirius-root/init

./initcorefs.sh

cd ./build/sirius-root/
find . -print0 | cpio --null -ov --format=newc > ../initrd.cpio
cd ../

gzip initrd.cpio
cd $base
