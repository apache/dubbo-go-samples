FROM golang:1.24-alpine AS builder
ENV TIMEZONE=Asia/Shanghai
WORKDIR /app

# Copy go mod files (build context is repo root)
COPY go.mod go.sum ./

# Copy only necessary source code (avoid copying deploy runtime data)
COPY apisix/go-server ./apisix/go-server/
COPY apisix/proto ./apisix/proto/

# Download dependencies
RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download

# Build
WORKDIR /app/apisix/go-server
RUN go build -o /app/server ./cmd/server.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server /app/server

EXPOSE 20001
ENTRYPOINT ["/app/server"]
