VERSION="v1.0.1"

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

build-web: install-dependencies
	cd web && \
	rm -rf dist && \
	npm run build

run: build-web
	go run ./dashboard --data-path ./data

static:
	CGO_ENABLED=0 go build -trimpath -ldflags="-X 'github.com/siacentral/sia-host-dashboard/dashboard/build.version=${VERSION}' -X 'github.com/siacentral/sia-host-dashboard/dashboard/build.gitRevision=${GIT_REVISION}' -X 'github.com/siacentral/sia-host-dashboard/dashboard/build.buildTime=${BUILD_TIME}' -s -w" -tags='netgo timetzdata' -o bin/ ./dashboard

release: build-web
	./release.sh