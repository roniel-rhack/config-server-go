FROM golang:1.25 AS builder
LABEL authors="roniel"

WORKDIR /app

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY ./src/ .

RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o configServer .

FROM alpine:3 AS production
LABEL authors="roniel"

RUN adduser -D -u 1000 appuser

WORKDIR /app

COPY --from=builder /app/configServer .
RUN chmod +x ./configServer

USER appuser

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
  CMD wget -qO- http://localhost:8888/health || exit 1

ENTRYPOINT ["./configServer"]
