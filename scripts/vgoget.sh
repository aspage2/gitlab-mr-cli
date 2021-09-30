#!/bin/sh

set -e
if ! [ $# -eq 2 ]; then
	echo 'usage: vgoget cmdpackage[@version] install_location'
	exit 2
fi

install_path="$(realpath "$2")"
d="$(mktemp -d)"
cd "$d" || exit
name="$(basename "$1")"
GOPATH="$d" go get "$1"
cp "$d/bin/$name" "$install_path"
