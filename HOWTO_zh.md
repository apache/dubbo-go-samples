# 如何运行

目前有三种方式来运行 dubbo-go 的示例:

1. 通过 makefile 快速开始: 在工程 "build" 子目录中提供了一个通用的 makefile。这个 makefile 可用于快速运行工程中的每一个示例。同时，由于这个 makefile 的存在，现在有机会可以把所有的示例串起来自动运行，从而使得通过该工程来做自动的 dubbo-go 的集成测试成为了可能。
2. 在 IDE 中快速开始，这也是**推荐**的方式: 在工程 ".run" 子目录下，提供了所有示例的 GoLand 运行配置文件，因此用户可以简单在 IDE 中单击运行所有的示例。
3. 在 IDE 中手工配置并运行: 为了完整性的目的，也为了万一您不使用 GoLand 而使用其他的 IDE，这里也提供了如何一步一步的配置的指南，帮助用户理解如何在 IDE 中配置，运行或者调试 dubbo-go 的示例。   

### 1. 通过 makefile 快速开始

*前置条件：需要 docker 环境就绪*

下面我们将使用 "attachment" 作为示例:

1. **启动注册中心（比如 zookeeper）**
   
   ```bash
   cd attachment/go-server
   make -f ../../build/Makefile docker-up 
   ```
   
   当看到类似下面的输出信息时，就表明 zookeeper server 启动就绪了。
   
   ```bash
   >  Starting dependency services with docker/docker-compose.yml
   Creating network "docker_default" with the default driver
   Creating docker_zookeeper_1 ... done
   ```
   
   如果要停掉注册中心，可以通过运行以下的命令完成:
   
   ```bash
   cd attachment/go-server
   make -f ../../build/Makefile docker-down
   ```
   
2. **启动服务提供方**
   
    ```bash
    cd attachment/go-server
    make -f ../../build/Makefile start
    ```
   
   当看到类似下面的输出信息时，就表明服务提供方启动就绪了。

   ```bash
   >  Buiding application binary: dist/darwin_amd64/release/go-server
   >  Starting application go-server, output is redirected to dist/darwin_amd64/release/go-server.log
     >  PID: 86428
   ```

   `go-server` 的输出信息可以在 'dist/darwin_amd64/release/go-server.log' 中找到。 
   
3. **运行服务调用方**
   
    ```bash
   cd attachment/go-client
   make -f ../../build/Makefile run 
   ```

   当以下的信息输出时，说明 `go-client` 调用 `go-server` 成功。

   ```bash
   >  Buiding application binary: dist/darwin_amd64/release/go-client
   >  Running application go-client, output is redirected to dist/darwin_amd64/release/go-client.log
   ...
   2020-10-27T14:51:37.520+0800    DEBUG   dubbo/dubbo_invoker.go:144      result.Err: <nil>, result.Rest: &{A001 Alex Stocks 18 2020-10-27 14:51:37.52 +0800 CST}
   2020-10-27T14:51:37.520+0800    DEBUG   proxy/proxy.go:177      [makeDubboCallProxy] result: &{A001 Alex Stocks 18 2020-10-27 14:51:37.52 +0800 CST}, err: <nil>
   response result: &{A001 Alex Stocks 18 2020-10-27 14:51:37.52 +0800 CST}
   ```
   
3. **集成测试**
   本项目 dubbo-go-samples 除了用来展示如何使用 dubbo-go 中的功能和特性之外，还被用于 apache/dubbo-go 的集成测试。可以按照以下的步骤来运行针对 `go-server` 设计的集成测试:

   ```bash
   cd attachment/go-server
   make -f ../../build/Makefile integration
   ```

   当以下信息输出时，说明集成测试通过。

   ```bash
   >  Running integration test for application go-server
   ...
   --- PASS: TestGetUser (0.00s)
   PASS
   ok      github.com/apache/dubbo-go-samples/attachment/go-server/tests/integration   3.603s
   ```
   
4. **关闭并清理**
   ```bash
   cd attachment/go-server
   make -f ../../build/Makefile clean docker-down
   ```

*以下的两种运行方式都与 IDE 有关。这里我们以 Intellij GoLand 为例来讨论。*

### 2. 在 IDE 中快速开始

一旦在 GoLand 中打开本项目，可以发现，在 "Run Configuration" 弹出菜单中已经存在了一系列事先配置好了的用来运行相关服务提供方和调用方的选项，例如："helloworld-go-server" 和 "helloworld-go-client"。

![run-configuration.png](.images/run-configurations.png)

可以选择其中的任意一个快速启动相关示例。当然在运行之前，假设需要的注册中心已经事先启动了，不然用例将会失败。您可以选择手动自行启动的方式，也可以利用工程中提供的 "docker-compose.yml" 在启动注册中心的 docker 实例。选择后者的话，可以参考[第三种方式](#3-manually-run-in-ide)中的细节。

### 3. 在 IDE 中手工运行

这里以 *Intellij GoLand* 为例。在 GoLand 中打开 dubbo-go-samples 工程之后，按照以下的步骤来运行/调试本示例:

1. **启动 zookeeper 服务器**

   打开 "attachment/go-server/docker/docker-compose.yaml" 这个文件，然后点击位于编辑器左边 gutter 栏位中的 ▶︎▶︎ 图标运行，"Service" Tab 应当会弹出并输出类似下面的文本信息:
   ```
   Deploying 'Compose: docker'...
   /usr/local/bin/docker-compose -f .../dubbo-go-samples/attachment/go-server/docker/docker-compose.yml up -d
   Creating network "docker_default" with the default driver
   Creating docker_zookeeper_1 ...
   'Compose: docker' has been deployed successfully.
   ```

2. **启动服务提供方**

   打开 "attachment/go-server/cmd/server.go" 文件，然后点击左边 gutter 栏位中紧挨着 "main" 函数的 ▶︎ 图标，并从弹出的菜单中选择 "Modify Run Configuration..."，并确保以下配置的准确:
   * Working Directory: "attachment/go-server" 目录的绝对路径，比如： */home/dubbo-go-samples/attachment/go-server*
   * Environment: CONF_PROVIDER_FILE_PATH=conf/server.yml, 另外也可以指定这个环境变量 "APP_LOG_CONF_FILE=conf/log.yml"

   这样示例中的服务端就准备就绪，随时可以运行了。

3. **运行服务消费方**

   打开 "attachment/go-client/cmd/client.go" 这个文件，然后从左边 gutter 栏位中点击紧挨着 "main" 函数的 ▶︎ 图标，然后从弹出的菜单中选择 "Modify Run Configuration..."，并确保以下配置的准确:
   * Working Directory: "attachment/go-client" 目录的绝对路径，比如： */home/dubbo-go-samples/attachment/go-client*
   * Environment: CONF_CONSUMER_FILE_PATH=conf/client.yml, 另外也可以指定这个环境变量 "APP_LOG_CONF_FILE=conf/log.yml"

   然后就可以运行并调用远端的服务了，如果调用成功，将会有以下的输出:
   ```
   [2021-02-03/16:19:30 main.main: client.go: 66] response result: &{A001 Alex Stocks 18 2020-02-04 16:19:30.422 +0800 CST}
   ```

如果需要调试该示例或者 dubbo-go 框架，可以在 IDE 中从 "Run" 切换到 "Debug"。如果要结束的话，直接点击 ◼︎ 就好了。

