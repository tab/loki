FROM golang:1.23-alpine as base-backend

ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=0

RUN apk add --no-cache --update git tzdata ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

RUN go build -o /app/loki /app/cmd/loki/main.go

CMD ["/app/loki"]
