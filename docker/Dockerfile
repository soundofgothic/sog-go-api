FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/ ./cmd/...

FROM alpine:3.18
RUN apk add --no-cache gettext

COPY --from=builder /app/sog-backend /app/
COPY docker/config-template.yaml /app/
COPY docker/entrypoint.sh /app/

ENTRYPOINT ["/app/entrypoint.sh"]
