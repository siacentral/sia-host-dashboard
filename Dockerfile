# build web
FROM node:13-alpine AS buildnode

WORKDIR /web

COPY ./web/package*.json ./

RUN npm install

COPY ./web .

RUN npm run build

# build daemon
FROM golang:alpine AS buildgo

WORKDIR /app

COPY . .
COPY --from=buildnode /web/dist web/dist

RUN apk -U --no-cache add git gcc make ca-certificates \
	&& update-ca-certificates \
	&& make static

# production
FROM scratch

COPY --from=buildgo /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=buildgo /app/bin/dashboard /

ENV SIA_API_ADDR="localhost:9980"

ENTRYPOINT [ "/dashboard", "--std-out", "--data-path", "/data" ]