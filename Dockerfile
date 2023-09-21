# syntax=docker/dockerfile:1.4.0
LABEL org.opencontainers.image.source="https://github.com/mgufrone/jenkins-bot-go" \
    org.opencontainers.image.licenses="MIT"
FROM golang:1.20 as build
WORKDIR /code
RUN --mount=type=bind,source=./go.mod,target=go.mod \
    --mount=type=bind,source=./go.sum,target=go.sum \
    --mount=type=bind,source=./pkg,target=pkg \
    --mount=type=bind,source=./src,target=src \
    --mount=type=bind,source=./main.go,target=main.go \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download && go build -o app .

FROM gcr.io/distroless/base
WORKDIR /code
COPY --from=build ./code/app /code/app
CMD ["/code/app"]
