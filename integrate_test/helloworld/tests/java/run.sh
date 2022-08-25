# !/bin/bash

if [ -z "$1" ]; then
  echo "Provide test directory please, like : ./run.sh $(pwd)/helloworld/java-server"
  exit
fi

#cd $1
#set -e
#
#mvn test