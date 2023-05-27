curl 127.0.0.1:2181
res=$?
passCode=52
while [ "$res" != "$passCode" ]; do
  sleep 5
  curl 127.0.0.1:2181
  res=$?
done

sleep 5
curl http://127.0.0.1:8848/nacos/v1/console/health/liveness
sleep 10

curl http://127.0.0.1:8090
