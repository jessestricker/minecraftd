#! /usr/bin/bash

set -e -u

arch=$(go env GOARCH)
package_dir=_package
package_file="minecraftd-$arch.tar.gz"

install -D minecraftd $package_dir/usr/local/bin/minecraftd >&2
install -D -m 644 minecraftd.service $package_dir/usr/local/lib/systemd/system/minecraftd.service >&2
install -D -m 644 minecraftd.toml $package_dir/usr/local/share/minecraftd/minecraftd.toml >&2
tar --create --gzip --owner 0 --group 0 --file "$package_file" --directory $package_dir . >&2
tar --list --verbose --file "$package_file" >&2

echo "$package_file"
