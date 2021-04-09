# before
# run:    kubectl apply -f ./zk.yml -n dubbo-workplace
kubectl create ns dubbo-workplace
kubectl apply -f ./deploy/crd.yml -n dubbo-workplace
kubectl apply -f ./deploy/zk.yml -n dubbo-workplace

GOOS=linux go build -o ./go-server/cmd/server ./go-server/cmd/
docker build --no-cache -t k8s-uniform-router-demo-server ./go-server
rm ./go-server/cmd/server
kubectl delete -f ./deploy/server.yml -n dubbo-workplace
kubectl apply -f ./deploy/server.yml -n dubbo-workplace

GOOS=linux go build -o ./go-server2/cmd/server ./go-server2/cmd/
docker build --no-cache -t k8s-uniform-router-demo-server2 ./go-server2
rm ./go-server2/cmd/server
kubectl delete -f ./deploy/server2.yml -n dubbo-workplace
kubectl apply -f ./deploy/server2.yml -n dubbo-workplace

GOOS=linux go build -o ./go-client/cmd/client ./go-client/cmd/
docker build --no-cache -t k8s-uniform-router-demo-client ./go-client
rm ./go-client/cmd/client
kubectl delete -f ./deploy/client.yml -n dubbo-workplace
kubectl apply -f ./deploy/client.yml -n dubbo-workplace
