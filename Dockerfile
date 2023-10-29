# syntax = docker/dockerfile:1.4
ARG GO_VERSION=1.21
ARG GOLANGCI_LINT_VERSION=v1.52
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS builder

ENV GO111MODULE=on \
    GOOS=linux

WORKDIR /build
RUN --mount=type=bind,src=go.mod,dst=go.mod \
    --mount=type=bind,src=go.sum,dst=go.sum \
    --mount=type=cache,dst=/go/pkg/mod \
    go mod download -x

FROM builder as bin
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=bind,dst=. \
    --mount=type=cache,dst=/go/pkg/mod  \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/main .


FROM scratch

WORKDIR /www

COPY ./database/ /www/database/
COPY ./public/ /www/public/
COPY ./storage/ /www/storage/
COPY ./resources/ /www/resources/
COPY --from=bin /bin/main /www/
#COPY --from=builder /build/.env /www/.env

ENTRYPOINT ["/www/main"]
