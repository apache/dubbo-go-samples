FROM golang:latest as build
ENV GO111MODULE on
ENV CONF_PROVIDER_FILE_PATH  /release/conf/server.yml
ADD ./go-server/dist/linux_amd64/ /
ENTRYPOINT exec /release/go-server