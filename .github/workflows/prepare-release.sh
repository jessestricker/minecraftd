#! /usr/bin/bash

set -e -u

version=$1

for package_file in minecraftd-*.tar.gz; do
    versioned_package_file=${package_file%.tar.gz}-$version.tar.gz
    mv "$package_file" "$versioned_package_file"
done
