P_DIR=$(pwd)/go-server
make GOOS="linux" PROJECT_DIR=$P_DIR PROJECT_NAME=$(basename $P_DIR) BASE_DIR=$P_DIR/dist -f ../build/Makefile build
docker build --no-cache -t dubbogo-docker-sample .
docker run --name zkserver -p 2181:2181 --restart always -d zookeeper:3.4.9
docker run -e DUBBO_IP_TO_REGISTRY=127.0.0.1  -p 20000:20000  --link zkserver:zkserver dubbogo-docker-sample
