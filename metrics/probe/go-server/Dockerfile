FROM alpine:3.19
WORKDIR /app
COPY metrics/probe/go-server/probeApp /app/server

EXPOSE 20000 22222
ENTRYPOINT ["/app/server"]
