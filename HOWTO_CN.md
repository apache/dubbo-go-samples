# 如何运行

[English](README.md) | 中文

目前有三种方式来运行 dubbo-go 的示例:

1. 通过 bash 命令快速开始: 通过简单的命令行启动样例以及进行单元测试
2. 在 IDE 中快速开始，这也是**推荐**的方式: 在工程 ".run" 子目录下，提供了所有示例的 GoLand 运行配置文件，因此用户可以简单在 IDE 中单击运行所有的示例。
3. 在 IDE 中手工配置并运行: 为了完整性的目的，也为了万一您不使用 GoLand 而使用其他的 IDE，这里也提供了如何一步一步的配置的指南，帮助用户理解如何在 IDE 中配置，运行或者调试 dubbo-go 的示例。   

### 1. 通过 命令行 快速开始

*前置条件：需要 docker 环境就绪*

下面我们将使用 "helloworld" 作为示例:

1. **启动注册中心（比如 zookeeper）**
   
   ```bash
   make -f Makefile docker-up 
   ```
   
   当看到类似下面的输出信息时，就表明 zookeeper server 启动就绪了。
   
   ```bash
   >  Starting dependency services with ./docker-compose.yml
   Docker Compose is now in the Docker CLI, try `docker compose up`
   
   Creating network "dockercompose_default" with the default driver
   Creating dockercompose_zookeeper_1 ... done
   Creating etcd                      ... done
   Creating nacos-standalone          ... done
   ```
   
   如果要停掉注册中心，可以通过运行以下的命令完成:
   
   ```bash
   make -f Makefile docker-down
   ```
   
2. **启动服务提供方**
   
    ```bash
    cd helloworld/go-server/cmd
    export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"
    go run .
    ```
   
   当看到类似下面的输出信息时，就表明服务提供方启动就绪了。

   ```bash
   2021/10/27 00:33:10 Connected to 127.0.0.1:2181
   2021/10/27 00:33:10 Authenticated: id=72057926938066944, timeout=10000
   2021/10/27 00:33:10 Re-submitting `0` credentials after reconnec
   ```
 
3. **运行服务调用方**
   
    ```bash
   cd helloworld/go-client/cmd
   export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"
   go run .
   ```

   当以下的信息输出时，说明 `go-client` 调用 `go-server` 成功。

   ```bash
   2021-10-27T00:40:44.879+0800    DEBUG   triple/dubbo3_client.go:106     TripleClient.Invoke: get reply = name:"Hello laurence" id:"12345" age:21 
   2021-10-27T00:40:44.879+0800    DEBUG   proxy/proxy.go:218      [makeDubboCallProxy] result: name:"Hello laurence" id:"12345" age:21 , err: <nil>
   2021-10-27T00:40:44.879+0800    INFO    cmd/client.go:51        client response result: name:"Hello laurence" id:"12345" age:21
   ```
   
4. **集成测试流程（当前 CI 逻辑）**
   现在的集成测试是脚本驱动的。单个 sample 可执行：
   ```bash
   ./integrate_test.sh <sample-path>
   ```
   例如：
   ```bash
   ./integrate_test.sh direct
   ```

   脚本会按以下顺序执行：
   1. 启动 `go-server`
   2. 运行 `go-client`
   3. 运行 `java-client`（若存在）
   4. 停止 `go-server`
   5. 启动 `java-server`（若存在）
   6. 运行 `java-client`
   7. 运行 `go-client`

   如果环境里没有 Maven（`mvn`），会自动跳过所有 Java 阶段，只运行 Go 阶段。

   若要本地执行完整 CI 样例列表：
   ```bash
   ./start_integrate_test.sh
   ```

   `start_integrate_test.sh` 会：
   - 通过仓库根目录 `docker-compose.yml` 启动依赖
   - 逐个执行 `./integrate_test.sh ...`
   - 在结束时（或失败时）回收依赖容器
   
7. **关闭并清理**
   ```bash
   make -f Makefile clean docker-down
   ```

*以下的两种运行方式都与 IDE 有关。这里我们以 Intellij GoLand 为例来讨论。*

### 2. 在 IDE 中快速开始

一旦在 GoLand 中打开本项目，可以发现，在 "Run Configuration" 弹出菜单中已经存在了一系列事先配置好了的用来运行相关服务提供方和调用方的选项，例如："helloworld-go-server" 和 "helloworld-go-client"。

![run-configuration.png](.images/run-configurations.png)

可以选择其中的任意一个快速启动相关示例。当然在运行之前，假设需要的注册中心已经事先启动了，不然用例将会失败。您可以选择手动自行启动的方式，也可以利用工程中提供的 "docker-compose.yml" 在启动注册中心的 docker 实例。选择后者的话，可以参考[第三种方式](#3-manually-run-in-ide)中的细节。

### 3. 在 IDE 中手工运行

这里以 *Intellij GoLand* 为例。在 GoLand 中打开 dubbo-go-samples 工程之后，按照以下的步骤来运行/调试本示例:

1. **启动 zookeeper 服务器**

   打开 "docker-compose.yml" 这个文件，然后点击位于编辑器左边 gutter 栏位中的 ▶︎▶︎ 图标运行，"Service" Tab 应当会弹出并输出类似下面的文本信息:
   ```
   Deploying 'Compose: docker'...
   /usr/local/bin/docker-compose -f ...docker-compose.yml up -d
   Creating network "docker_default" with the default driver
   Creating docker_zookeeper_1 ...
   'Compose: docker' has been deployed successfully.
   ```

2. **启动服务提供方**

   打开 "helloworld/go-server/cmd/server.go" 文件，然后点击左边 gutter 栏位中紧挨着 "main" 函数的 ▶︎ 图标，并从弹出的菜单中选择 "Modify Run Configuration..."，并确保以下配置的准确:
   * Working Directory: "helloworld/go-server" 目录的绝对路径，比如： */home/dubbo-go-samples/helloworld/go-server*
   * Environment: DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"

   这样示例中的服务端就准备就绪，随时可以运行了。

3. **运行服务消费方**

   打开 "helloworld/go-client/cmd/client.go" 这个文件，然后从左边 gutter 栏位中点击紧挨着 "main" 函数的 ▶︎ 图标，然后从弹出的菜单中选择 "Modify Run Configuration..."，并确保以下配置的准确:
   * Working Directory: "helloworld/go-client" 目录的绝对路径，比如： */home/dubbo-go-samples/helloworld/go-client*
   * Environment: DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"

   然后就可以运行并调用远端的服务了，如果调用成功，将会有以下的输出:
   ```
   [2021-02-03/16:19:30 main.main: client.go: 66] response result: &{A001 Alex Stocks 18 2020-02-04 16:19:30.422 +0800 CST}
   ```

如果需要调试该示例或者 dubbo-go 框架，可以在 IDE 中从 "Run" 切换到 "Debug"。如果要结束的话，直接点击 ◼︎ 就好了。
