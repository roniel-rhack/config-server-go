FROM golang:1.22.2 as builder
LABEL authors="roniel"

WORKDIR /app

COPY ./src/ .

RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o configServer .

FROM alpine:3 as production
LABEL authors="roniel"

WORKDIR /app

COPY --from=builder /app/configServer .
COPY --from=builder /app/application.yaml .
RUN chmod +x ./configServer && chmod 777 ./hgw-stddev.yaml

ENTRYPOINT ["./configServer"]