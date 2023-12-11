FROM golang:1.21 AS builder

WORKDIR /pow
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

FROM alpine:latest

COPY --from=builder /pow/server /server
COPY --from=builder /pow/config/server.yml /etc/pow/server.yml

ENV POW_CONFIG_FILE server.yml
ENV POW_CONFIG_PATH /etc/pow

EXPOSE 8234

ENTRYPOINT /server