FROM golang:latest@sha256:43813e40aba2a453d4fac2ab5f6357273b6491d5225f8b703eb1a9c70aa4055a as build
WORKDIR /
COPY ./cmd/server /
COPY ./conf /conf
ENV GO111MODULE on
ENV CONF_PROVIDER_FILE_PATH  /conf/server.yml
ENV APP_LOG_CONF_FILE  /conf/log.yml
ENTRYPOINT exec ./server
