curl 127.0.0.1:2379
res=$?
passCode=52
while [ "$res" != "$passCode" ];do
  sleep 5
  curl 127.0.0.1:2379
  res=$?
done

sleep 5
