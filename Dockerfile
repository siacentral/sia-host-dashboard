# build web app
FROM node:12 AS buildnode

WORKDIR /web

COPY ./web/package*.json ./

RUN npm install

COPY ./web .

RUN npm run build

# build wasm
FROM golang:1.13-alpine AS buildgo

WORKDIR /daemon

COPY ./daemon .
COPY --from=buildnode /web/dist /daemon/dist

RUN go get
RUN go run generate/assets_generate.go /daemon/dist
RUN go build -o ./release/dashboard main.go

# production
FROM alpine:latest

COPY --from=buildgo /daemon/release/dashboard /usr/bin/dashboard

ENV SIA_API_ADDR="localhost:9980"

ENTRYPOINT /usr/bin/dashboard --std-out --data-path "/data" --sia-api-addr $SIA_API_ADDR