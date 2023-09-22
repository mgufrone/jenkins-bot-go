# syntax=docker/dockerfile:1.4.0
FROM golang:1.20 as build
WORKDIR /code
RUN --mount=type=bind,source=./go.mod,target=go.mod \
    --mount=type=bind,source=./go.sum,target=go.sum \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download
RUN --mount=type=bind,source=.,target=. \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o /bin/app .

FROM gcr.io/distroless/base
LABEL org.opencontainers.image.source="https://github.com/mgufrone/jenkins-bot-go" \
    org.opencontainers.image.licenses="MIT"
WORKDIR /code
COPY --from=build /bin/app /code/app
CMD ["/code/app"]
