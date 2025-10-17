# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ---- Build stage ----
FROM golang:1.23.0 AS builder
WORKDIR /app

# Copy go.mod first to leverage Docker build cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pgw-cleaner ./cmd/pgw-cleaner

# ---- Runtime stage ----
FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata

# Copy binary
COPY --from=builder /app/pgw-cleaner /bin/pgw-cleaner

# Non-root user
USER 1000

ENTRYPOINT ["/bin/pgw-cleaner"]
