# syntax=docker/dockerfile:1
FROM golang:1.18.0-alpine3.15 as builder
ENV TIMEZONE Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
&& apk update && apk add git 
WORKDIR /aixichen/helloworld
RUN mkdir -p /aixichen/helloworld
#添加文件
COPY conf /aixichen/helloworld
##编译
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct \
&& cd /aixichen/helloworld/ && go build -o aixichen-service

FROM alpine:3.15
RUN mkdir -p /aixichen/helloworld
COPY --from=builder /aixichen/helloworld/aixichen-service /aixichen/helloworld
WORKDIR /aixichen/helloworld

EXPOSE 80
ENTRYPOINT ["./aixichen-service"]
