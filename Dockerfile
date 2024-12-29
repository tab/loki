FROM golang:1.23-alpine3.21 AS builder

ENV CGO_ENABLED=0

RUN apk add --no-cache --update git tzdata ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

COPY . ./
RUN go build -o /app/loki /app/cmd/loki/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/loki /app/loki

CMD ["/app/loki"]
