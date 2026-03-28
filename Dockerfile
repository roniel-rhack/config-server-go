FROM golang:1.23 AS builder
LABEL authors="roniel"

WORKDIR /app

COPY ./src/ .

RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o configServer .

FROM alpine:3 AS production
LABEL authors="roniel"

WORKDIR /app

COPY --from=builder /app/configServer .
RUN chmod +x ./configServer

ENTRYPOINT ["./configServer"]
