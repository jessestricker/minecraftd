#! /usr/bin/bash

set -e -u

ARCHITECTURES=(amd64 arm64)

version=$1

# re-create debian/changelog
rm -f debian/changelog
EDITOR=true debchange --newversion "$version" --create --empty --package minecraftd --controlmaint ' '

# build packages
for arch in "${ARCHITECTURES[@]}"; do
    dpkg-buildpackage --build=binary --host-arch "$arch" --no-sign
done
