BUILD_TIME=$(shell date)
GIT_REVISION=$(shell git rev-parse --short HEAD)

all: release

install-dependencies: 
	cd web && \
	rm -rf node_modules && \
	npm i

lint-web:
	cd web && \
	npm run lint -- --fix

lint-daemon:
	cd daemon && \
	go get golang.org/x/lint/golint && \
	golint -min_confidence=1.0 -set_exit_status $(pkgs)

lint: lint-web lint-daemon

pack:
	cd web && \
	rm -rf node_modules dist && \
	npm i && \
	npm run build && \
	cd ../daemon && \
	go run generate/assets_generate.go ../web/dist

run: install-dependencies lint-web lint-daemon pack
	cd daemon && \
	go run main.go

build: install-dependencies lint-web lint-daemon pack build
	./release.sh $(GOHOSTOS) $(GOHOSTARCH)

release: install-dependencies lint-web lint-daemon pack build
	./release.sh
