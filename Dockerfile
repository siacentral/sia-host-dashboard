# build web app
FROM node:13-alpine AS buildnode

WORKDIR /web

COPY ./web/package*.json ./

RUN npm install

COPY ./web .

RUN npm run build

# build daemon
FROM golang:1.13-alpine AS buildgo

RUN apk update && apk upgrade && apk add --no-cache alpine-sdk

WORKDIR /app

COPY . .
COPY --from=buildnode /web/dist ./dist

RUN go run generate/assets_generate.go ./dist
RUN go build -trimpath -o ./release/dashboard \
	-ldflags="-X 'github.com/siacentral/host-dashboard/daemon/build.GitRevision=`git rev-parse --short HEAD`' -X 'github.com/siacentral/host-dashboard/daemon/build.BuildTimestamp=`git show -s --format=%ci HEAD`'" \
	./daemon

# production
FROM alpine:latest

COPY --from=buildgo /app/release/dashboard /usr/bin/dashboard

ENV SIA_API_ADDR="localhost:9980"

ENTRYPOINT /usr/bin/dashboard --std-out --data-path "/data" --sia-api-addr $SIA_API_ADDR