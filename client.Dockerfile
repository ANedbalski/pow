FROM golang:1.21 AS builder

WORKDIR /pow
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o client ./cmd/client/main.go

FROM alpine:latest

COPY --from=builder /pow/client /client
COPY --from=builder /pow/config/client.yml /etc/pow/client.yml

ENV POW_CONFIG_FILE client.yml
ENV POW_CONFIG_PATH /etc/pow

ENTRYPOINT /client