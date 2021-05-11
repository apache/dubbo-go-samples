FROM golang:latest as build
WORKDIR /
COPY ./cmd/server /
COPY ./conf /conf
ENV GO111MODULE on
ENV CONF_PROVIDER_FILE_PATH  /conf/server.yml
ENV APP_LOG_CONF_FILE  /conf/log.yml
ENTRYPOINT exec ./server
