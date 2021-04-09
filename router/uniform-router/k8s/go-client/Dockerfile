FROM golang:latest as build
WORKDIR /
COPY ./cmd/client /
COPY ./conf /conf
ENV GO111MODULE on
ENV APP_LOG_CONF_FILE  /conf/log.yml
ENV CONF_CONSUMER_FILE_PATH  /conf/client.yml
ENV CONF_VIRTUAL_SERVICE_FILE_PATH /conf/virtual_service.yml
ENV CONF_DEST_RULE_FILE_PATH /conf/dest_rule.yml
ENTRYPOINT exec ./client