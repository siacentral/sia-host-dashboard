#!/bin/bash

set -e

buildtime=$(git show -s --format=%ci HEAD)
gitrevision=$(git rev-parse --short HEAD)
buildflags="-X 'github.com/siacentral/sia-host-dashboard/daemon/build.GitRevision=$gitrevision' -X 'github.com/siacentral/sia-host-dashboard/daemon/build.BuildTimestamp=$buildtime'"

build() {
	os=$1
	arch=$2
	folder=./releases/dashboard-$os-$arch
	bin=dashboard

	if [ "$os" == "windows" ]; then
		bin=dashboard.exe
	fi

	echo "Building $os $arch..."

	rm -rf $folder
	mkdir -p $folder
	GOOS=${os} GOARCH=${arch} go build -ldflags="$buildflags" -a -trimpath -o $folder/$bin ./daemon/daemon.go
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

	if [ "$os" == "linux" ]; then
		build linux arm64
	fi
done