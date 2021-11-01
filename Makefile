-include .env.local

VERSION="v1.0.1"

ifeq "$(shell git status --porcelain=v1 2>/dev/null)" "" 
GIT_REVISION=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell git show -s --format=%ci HEAD)
else
GIT_REVISION="$(shell git rev-parse --short HEAD)-devel"
BUILD_TIME=$(shell date)
endif

# get the current OS and arch
OS=$(shell GOOS= go env GOOS)
ARCH=$(shell GOARCH= go env GOARCH)

# get the target OS and arch
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

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
	CGO_ENABLED=0 go build -trimpath -ldflags="-X 'github.com/siacentral/sia-host-dashboard/dashboard/build.version=${VERSION}' -X 'github.com/siacentral/sia-host-dashboard/dashboard/build.gitRevision=${GIT_REVISION}' -X 'github.com/siacentral/sia-host-dashboard/dashboard/build.buildTime=${BUILD_TIME}' -s -w" -tags='netgo timetzdata' -o ./bin/ ./dashboard

package: static
	mkdir -p dist
	@if [ "${GOOS}" = "darwin" ] && [ "${OS}" != "darwin" ]; then \
		echo "darwin must be packaged on macOS"; \
	elif [ "${GOOS}" = "darwin" ] && [ "${OS}" = "darwin" ]; then \
		codesign --deep -s $(APPLE_CERT_ID) -o runtime bin/dashboard; \
		ditto -ck bin dist/host_dashboard_$(VERSION)_$(GOOS)_$(GOARCH).zip; \
		xcrun altool --notarize-app --primary-bundle-id com.siacentral.sia-host-dashboard --apiKey $(APPLE_API_KEY) --apiIssuer $(APPLE_API_ISSUER) --file dist/host_dashboard_$(VERSION)_$(GOOS)_$(GOARCH).zip; \
		rm -rf bin; \
	else \
		zip -qj dist/host_dashboard_$(VERSION)_$(GOOS)_$(GOARCH).zip bin/*; \
		rm -rf bin; \
	fi

# must be run on macOS
apple-notarize-status:
	@if [ "${OS}" = "darwin" ]; then \
		xcrun altool --notarization-info $(APPLE_NOTARIZE_UUID) --apiKey $(APPLE_API_KEY) --apiIssuer $(APPLE_API_ISSUER); \
	else \
		echo "Notarization is only supported on macOS"; \
	fi

release: build-web
	rm -rf bin dist
	@for OS in linux darwin windows; do \
		for ARCH in amd64 arm64; do \
			GOOS=$$OS GOARCH=$$ARCH make package; \
		done; \
	done
	rm -rf bin