sleep 10
curl http://127.0.0.1:8848/nacos/v1/console/health/liveness
curl 127.0.0.1:8848
res=$?
passCode=52
while [ "$res" != "$passCode" ];do
  sleep 5
  curl 127.0.0.1:8848
  res=$?
done