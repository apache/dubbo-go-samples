# Seata-go tcc example

## How to run

1. Start the seata-server and zookeeper service with the docker file

   ~~~shell
   cd seata-go/tcc
   docker-compose -f docker-compose.yml up -d seata-server zookeeper
   ~~~

2. Just execute the main function under tcc/client/cmd and tcc/server/cmd directory