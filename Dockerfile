# Multi-stage build for the Go API

FROM golang:1.22-alpine AS builder
WORKDIR /src

# Install CA certs for any HTTPS downloads during build
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/main.go


# Runtime image
FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 app

COPY --from=builder /out/app /app/app
# Copy default config; compose sets CONFIG_PATH to point to docker config
COPY config.yaml /app/config.yaml

EXPOSE 5000
USER app

ENTRYPOINT ["/app/app"]
