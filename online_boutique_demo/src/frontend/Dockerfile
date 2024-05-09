FROM golang:1.19.0-alpine AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/hipstershop

# Install dependencies
RUN apk --update --no-cache add ca-certificates make protoc

# Build Go binary
COPY Makefile go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.io,direct/
RUN make init && go mod download
COPY . .
RUN make proto tidy

# Skaffold passes in debug-oriented compiler flags
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /go/src/hipstershop/frontend . 

# Deployment container
FROM scratch

COPY --from=builder /go/src/hipstershop/frontend /hipstershop/frontend
COPY ./templates ./templates
COPY ./static ./static

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

ENTRYPOINT ["/hipstershop/frontend"]