#!/bin/bash

set -e

build() {
	os=$1
	arch=$2
	folder=./releases/dashboard-$os-$arch
	archive=./releases/dashboard-$os-$arch.zip
	bin=dashboard

	if [ "$os" == "windows" ]; then
		bin=dashboard.exe
	fi

	echo "Building $os $arch..."

	rm -rf $folder $archive
	mkdir -p $folder
	GOOS=${os} GOARCH=${arch} make static
	mv bin/$bin $folder/$bin
	zip $archive -j $folder/$bin
	rm -rf $folder
}

sys=(darwin linux windows)

if [ "$1" != "" ] && [ "$2" != "" ]; then
	build $1 $2
	exit 0
elif [ "$1" != "" ]; then
	sys=( $1 )
fi

for os in ${sys[@]}; do
	build $os amd64

	if [ "$os" = "linux" ] || [ "$os" = "darwin" ]; then
		build $os arm64
	fi
done