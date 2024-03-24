# Seata-go tcc example

## How to run?

1. Start the seata-server with the docker file.

   ```shell
   cd ../dockercompose
   docker-compose -f docker-compose.yml up -d seata-server
   ```

2. Just execute the main function under non-idl/client/cmd and non-idl/server/cmd directory.
