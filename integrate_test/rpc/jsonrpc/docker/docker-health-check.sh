echo "check zk:2181"
curl 127.0.0.1:2181
res=$?
passCode=52
while [ "$res" != "$passCode" ];do
  sleep 5
  curl 127.0.0.1:2181
  res=$?
done

echo "check zk:2182"
curl 127.0.0.1:2182
res=$?
passCode=52
while [ "$res" != "$passCode" ];do
  sleep 5
  curl 127.0.0.1:2182
  res=$?
done

#sleep 5