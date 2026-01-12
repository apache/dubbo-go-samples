FROM golang:1.24-alpine AS builder
ENV TIMEZONE Asia/Shanghai
WORKDIR /app

# Copy go mod files (build context is repo root)
COPY go.mod go.sum ./

# Copy mesh source code
COPY mesh ./mesh/

# Download dependencies
RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download

# Build
WORKDIR /app/mesh/go-client
RUN go build -o /app/client ./cmd/client.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/client /app/client

EXPOSE 50052
ENTRYPOINT ["/app/client"]

