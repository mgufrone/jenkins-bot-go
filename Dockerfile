FROM golang:1.16-alpine as builder

WORKDIR /app
ADD . /app
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /app/app /app/app
CMD ["/app/app"]