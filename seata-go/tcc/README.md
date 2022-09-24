# Seata-go tcc example

## How to run

1. Start the seata-server service with the docker file under the sample/dockercomposer folder

   ~~~shell
   cd seata-go/tcc
   docker-compose -f docker-compose.yml up -d seata-server zookeeper
   ~~~

2. Just execute the main function under samples/ in the root directory