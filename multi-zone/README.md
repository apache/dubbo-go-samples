# Run Steps

1. Start Zookeeper in docker environment.

    ```bash
    make -f ../build/Makefile docker-up
    ```
   
2. Start Server

    Start Hangzhou zone server.
    ```bash
   cd go-server-hz
   make -f ../../build/Makefile clean start  
   ```
   
   Start Shanghai zone server.
   ```bash
   cd go-server-sh
   make -f ../../build/Makefile clean start  
  
3. Run Consumer

    ```bash
    cd go-client
    make -f ../../build/Makefile run
    ```
   
4. Cleanup

   ```bash
   cd go-server-hz && \
        make -f ../../build/Makefile clean && \
        cd ..
   cd go-server-sh && \
        make -f ../../build/Makefile clean && \
        cd ..
   cd go-client && \
        make -f ../../build/Makefile clean && \
        cd ..
   make -f ../build/Makefile docker-down
   ```