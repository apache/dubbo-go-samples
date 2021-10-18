# !/bin/bash

if [ -z "$1" ]; then
  echo "Provide log file please, like : ./waiting_launch.sh log.sh"
  exit
fi

PREV_MD5=""
MD5=""

# wait 150s at most
for ((i=1; i<=15; i++));
do
  sleep 10s
  MD5=$(md5sum $1 | cut -d ' ' -f1)
  if [ "$PREV_MD5" = "$MD5" ]; then
    exit
  fi
  echo "waiting... log file md5: $MD5"
  PREV_MD5=$MD5
done

echo "java-server is not launched properly, the launching log will be outputted at below: "
cat $1