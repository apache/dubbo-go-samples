FROM golang:1.18.0-alpine3.15 as builder
ENV TIMEZONE Asia/Shanghai
WORKDIR /mesh/
#添加文件
COPY go.mod ./
COPY go.sum ./
COPY conf ./conf
COPY cmd ./cmd
##编译
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go mod download
RUN go build -o client-demo ./...

FROM alpine:3.15
WORKDIR /mesh
COPY --from=builder /mesh/client-demo ./client-demo
COPY --from=builder /mesh/conf .

ENV DUBBO_GO_CONFIG_PATH="./dubbogo.yml"
ENTRYPOINT ["./client-demo"]
