FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/ ./cmd/...

FROM alpine:3.21
RUN apk add --no-cache gettext

COPY --from=builder /app/sog-backend /app/
COPY docker/config-template.yaml /app/
COPY docker/entrypoint.sh /app/

ENTRYPOINT ["/app/entrypoint.sh"]
