FROM golang:1.19 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg
COPY openapi.yaml oapi.yaml Makefile ./
RUN apt-get install -y make
RUN make cloud-provider-manager lb-api

FROM debian:stable-slim as cloud-provider-manager
WORKDIR /

COPY --from=build /app/cloud-provider-manager /cloud-provider-manager

USER nobody:nogroup

ENTRYPOINT ["/cloud-provider-manager"]
