# How To Run

There are three ways to run dubbo-go samples:

1. Quick start with bash command: start the sample and perform unit testing through a simple command line
2. Quick start in IDE (**Recommended**): In ".run" subdirectory a couple of GoLand run configuration files are provided so that user can run each sample with just one click.
3. Manually config and run in IDE: For completeness purpose, a step-by-step instruction is also provided so that user can understand how to configure and run or debug a sample in IDE. 

## 1. Quick start with makefile

*Prerequisite: docker environment is required*

Here we use "helloworld" as an example:

1. **Get the root path of dubbo-go-samples**

   ```bash
   cd <PATH OF dubbo-go-samples>
   export DUBBO_GO_SAMPLES_ROOT_PATH=$(pwd)
   ```

2. **Start register server (e.g. zookeeper)**
   
   ```bash
   make -f build/Makefile docker-up 
   ```
   
   Once the following messages outputs, the zookeeper server is ready.
   
   ```bash
   >  Starting dependency services with ./integrate_test/dockercompose/docker-compose.yml
   Docker Compose is now in the Docker CLI, try `docker compose up`
   
   Creating network "dockercompose_default" with the default driver
   Creating dockercompose_zookeeper_1 ... done
   Creating etcd                      ... done
   Creating nacos-standalone          ... done
   ```
   
   To shut it down, simple run
   
   ```bash
   make -f build/Makefile docker-down
   ```
   
3. **Start server**
   
    ```bash
    cd helloworld/go-server/cmd
    export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"
    go run .
    ```
   
   Once the following messages outputs, the server is ready.

   ```bash
   2021/10/27 00:33:10 Connected to 127.0.0.1:2181
   2021/10/27 00:33:10 Authenticated: id=72057926938066944, timeout=10000
   2021/10/27 00:33:10 Re-submitting `0` credentials after reconnec
   ```

   The output of `go-server` can be found from 'dist/darwin_amd64/release/go-server.log'.
   
4. **Run client**
   
    ```bash
   cd helloworld/go-client/cmd
   export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"
   go run .
   ```

   Once the following messages outputs, the `go-client` calls the `go-server` successfully.

   ```bash
   2021-10-27T00:40:44.879+0800    DEBUG   triple/dubbo3_client.go:106     TripleClient.Invoke: get reply = name:"Hello laurence" id:"12345" age:21 
   2021-10-27T00:40:44.879+0800    DEBUG   proxy/proxy.go:218      [makeDubboCallProxy] result: name:"Hello laurence" id:"12345" age:21 , err: <nil>
   2021-10-27T00:40:44.879+0800    INFO    cmd/client.go:51        client response result: name:"Hello laurence" id:"12345" age:21
   ```
   
5. **Integration test**
   dubbo-go-samples is designed to serve the purposes of not only the showcases of how to use apache/dubbo-go but also the integration-test for apache/dubbo-go. To run integration test for `go-server`, run the following commands:

   Start the server first
   ```bash
   cd helloworld/go-server/cmd
   export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"
   go run .
   ```

   Then switch to the single test directory, set the environment variables, and then execute the single test
   ```bash
   cd integrate_test/helloworld/tests/integration
   export DUBBO_GO_CONFIG_PATH="../../../../helloworld/go-client/conf/dubbogo.yml"
   go test -v
   ```

   Once the following messages outputs, the integration tests pass.

   ```bash
   >  Running integration test for application go-server
   ...
   --- PASS: TestSayHello (0.01s)
   PASS
   ok      github.com/apache/dubbo-go-samples/integrate_test/helloworld/tests/integration  0.119s
   ```
   
7. **Shutdown and cleanup**
   ```bash
   make -f build/Makefile clean docker-down
   ```

*The following two ways are all relevant to IDE. Intellij GoLand is discussed here as an example.*

## 2. Quick start in IDE

Once open this project in GoLand, a list of pre-configured configures for both server and client can be found from "Run Configuration" pop up menu, for example: "helloworld-go-server" and "helloworld-go-client". 

![run-configuration.png](.images/run-configurations.png)

Feel free to pick one among them to run instantly. Of course a service register server is required otherwise the samples cannot run successfully. You may consider to either manually boot up the required register server, or use the provided "docker-compose.yml" to boot it in docker environment, from the following [section](#3-manually-run-in-ide) where you can find more details.

## 3. Manually run in IDE

After open dubbo-go-samples in GoLand, follow the steps below to run/debug this
example:

1. **Start up zookeeper server**

   Open "integrate_test/dockercompose/docker-compose.yml", and click ▶︎▶︎ icon in the gutter on the left side of the
   editor, then "Services" tab should pop up and shows the similar message below:
   ```
   Deploying 'Compose: docker'...
   /usr/local/bin/docker-compose -f .../dubbo-go-samples/helloworld/go-server/docker/docker-compose.yml up -d
   Creating network "docker_default" with the default driver
   Creating docker_zookeeper_1 ...
   'Compose: docker' has been deployed successfully.
   ```

2. **Start up service provider**

   Open "helloworld/go-server/cmd/server.go", and click ▶︎ icon just besides "main" function in the gutter on the left
   side, and select "Modify Run Configuration..." from the pop-up menu. Then make sure the following configs configured
   correctly:
    * Working Directory: the absolute path to "helloworld/go-server", for examples: *
      /home/dubbo-go-samples/helloworld/go-server*
    * Environment: DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"

   Then the sample server is ready to run.

3. **Run service consumer**

   Open "helloworld/go-client/cmd/client.go", and click ▶︎ icon just besides "main" function in the gutter on the left
   side, and select "Modify Run Configuration..." from the pop-up menu. Then make sure the following configs configured
   correctly:
    * Working Directory: the absolute path to "helloworld/go-client", for examples: *
      /home/dubbo-go-samples/helloworld/go-client*
    * Environment: DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"

   Then run it to call the remote service, you will observe the following message output:
   ```
   [2021-02-03/16:19:30 main.main: client.go: 66] response result: &{A001 Alex Stocks 18 2020-02-04 16:19:30.422 +0800 CST}
   ```

If you need to debug either the samples or dubbo-go, you may consider switch to **Debug** instead of **Run** in GoLand. To stop, simply click ◼︎ to shutdown everything.

