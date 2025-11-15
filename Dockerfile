# BUILD STAGE
FROM golang:1.25.4-trixie AS build

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download

COPY ./ .
RUN go build -o labwrangler main.go

# DEPLOY STAGE
FROM debian:13-slim AS deploy
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates

RUN useradd -rm -s /bin/sh -u 1000 labwrangler
RUN mkdir -p /etc/labwrangler && chown -R 1000:1000 /etc/labwrangler
USER labwrangler

COPY --from=build /usr/src/app/labwrangler /usr/bin

VOLUME ["/etc/labwrangler"]
ENTRYPOINT ["/usr/bin/labwrangler"]
