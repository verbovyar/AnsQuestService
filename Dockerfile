FROM golang:1.24.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o qa ./cmd/app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/qa /app/qa

COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY migrations ./migrations

COPY entrypoint.sh .

RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
