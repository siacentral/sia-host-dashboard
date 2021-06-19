ifeq "$(shell git status --porcelain=v1 2>/dev/null)" "" 
GIT_REVISION=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell git show -s --format=%ci HEAD)
else
GIT_REVISION="$(shell git rev-parse --short HEAD)-devel"
BUILD_TIME=$(shell date)
endif

all: release

lint-web:
	cd web && \
	npm run lint -- --fix

lint-daemon:
	golangci-lint run

lint: lint-web lint-daemon

install-dependencies:
	cd web && \
	npm i

build-web: install-depencencies
	cd web && \
	rm -rf dist && \
	npm run build

run: build-web
	go run daemon/daemon.go --data-path $(PWD)/data

static: 
	CGO_ENABLED=0 go build -trimpath -ldflags="-X 'github.com/siacentral/sia-host-dashboard/build.gitRevision=${GIT_REVISION}' -X 'github.com/siacentral/sia-host-dashboard/build.buildTime=${BUILD_TIME}' -buildid='' -s -w -extldflags '-static'" -tags='netgo timetzdata'  -o ./bin/dashboard ./daemon

release: lint-web lint-daemon build-web
	./release.sh