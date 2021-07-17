curl 127.0.0.1:2182
res=$?
passCode=52
while [ "$res" != "$passCode" ];do
  sleep 5
  curl 127.0.0.1:2182
  res=$?
done

curl 127.0.0.1:2183
res=$?
while [ "$res" != "$passCode" ];do
  sleep 5
  curl 127.0.0.1:2183
  res=$?
done

sleep 5
