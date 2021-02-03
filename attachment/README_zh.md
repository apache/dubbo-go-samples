# Attachment 示例

### 背景

可以通过 attachment 把用户的数据从 Dubbo 的客户端传递给服务端。在 Dubbo-go 中，attachment 在二者之间的传递是通过 `context.Context` 的机制来完成的。如果要使用 attachment，那么应当把 `context.Context` 作为要调用的服务方法的第一个参数，例如:

```go
GetUser func(ctx context.Context, req []interface{}, rsp *User) error
```

为了在客户端传递用户数据到服务端，首先需要在方法的第一个参数 `context.Context` 中放入一个 `map[string]interface{}` 的数据类型，该数据的 Key 值约定为 "attachment"。以下的代码片段展示了如何把用户自定义的的一个时间戳通过该 map 中的 "timestamp" 放入 attachment 中:

```go
ctx := context.WithValue(context.Background(), constant.AttachmentKey, 
	map[string]interface{}{"timestamp": time.Now()})
err := userProvider.GetUser(ctx, []interface{}{"A001"}, user)
```

在服务提供方，方法的第一个参数 `context.Context` 传入时携带了用户在客户端放入的自定义数据。下面的代码就是本例子中的服务端实现，主要展示了如何从 attachment 中提取用户自定义的数据:


```go
func (u *UserProvider) GetUser(ctx context.Context, req []interface{}) (*User, error) {
	t := time.Now()
	attachment := ctx.Value(constant.AttachmentKey).(map[string]interface{})
	if v, ok := attachment["timestamp"]; ok {
		t = v.(time.Time).Add(-1 * 365 * 24 * time.Hour)
	}

	rsp := User{"A001", "Alex Stocks", 18, t}
	return &rsp, nil
}
```

### 如何运行

#### 通过 Makefile 的方式快速开始

*先决条件：需要提前安装好 Docker*

1. **启动服务端**
    ```bash
    cd attachment/server
    make -f ../../build/Makefile docker-up start
    ```
2. **运行客户端**
    ```bash
   cd attachment/client
   make -f ../../build/Makefile run 
   ```
3. **集成测试**
   ```bash
   cd attachment/server
   make -f ../../build/Makefile integration
   ```
4. **清理现场**
   ```bash
   cd attachment/server
   make -f ../../build/Makefile clean docker-down
   ```

#### 在 IDE 中运行

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

