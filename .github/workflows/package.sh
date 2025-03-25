#! /usr/bin/bash

set -e -u

arch=$(go env GOARCH)
package_dir=_package
package_file="minecraftd-$arch.tar.gz"

install -D        minecraftd        $package_dir/usr/local/bin/minecraftd
install -D -m 644 minecraft.service $package_dir/usr/local/lib/systemd/system/minecraft.service
tar --create --gzip --owner 0 --group 0 --file "$package_file" --directory $package_dir .
tar --list --verbose --file "$package_file"

echo "$package_file"
