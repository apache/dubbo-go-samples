FROM golang:latest

ARG APP=go-server-order
ENV APP ${APP}
ENV GOPROXY https://goproxy.io,direct
## provider config
ENV CONF_PROVIDER_FILE_PATH /build/$APP/conf/server.yml
# cunsumer config
ENV CONF_CONSUMER_FILE_PATH /build/$APP/conf/client.yml
## log config
ENV APP_LOG_CONF_FILE /build/$APP/conf/log.yml
## seata config
ENV SEATA_CONF_FILE /build/$APP/conf/seata.yml

COPY . /build

RUN echo ${APP} \
    && go env -w GOSUMDB=off GO111MODULE=auto
RUN cd /build/$APP/cmd \
    && go build

ENTRYPOINT /build/$APP/cmd/cmd