FROM gcr.io/distroless/base
WORKDIR /app
COPY ./app /app/app
CMD ["/app/app"]