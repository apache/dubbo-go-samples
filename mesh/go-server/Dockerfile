FROM golang:1.18.0-alpine3.15 as builder
ENV TIMEZONE Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
&& apk update && apk add git 
WORKDIR /mesh
#添加文件
COPY conf ./
##编译
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct \
&& cd /mesh/ && go build -o server-demo

FROM alpine:3.15
WORKDIR /mesh
COPY --from=builder /mesh/server-demo ./server-demo

EXPOSE 20000
ENTRYPOINT ["./server-demo"]
