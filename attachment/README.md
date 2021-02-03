# Attachment Example

### Background

A Dubbo client can pass user data to the remote Dubbo server via attachment. Dubbo-go leverages `context.Context` as the attachment between consumer and provider. In order to use attachment, a `context.Context` is required to introduce as the first parameter of the service method, for example:

```go
GetUser func(ctx context.Context, req []interface{}, rsp *User) error
```

To pass the user data from the client side, a `map[string]interface{}` should be put into the `context.Context` with the key "attachment". The following code snippet shows how a user data "timestamp" is put into the attachment:

```go
ctx := context.WithValue(context.Background(), constant.AttachmentKey, 
	map[string]interface{}{"timestamp": time.Now()})
err := userProvider.GetUser(ctx, []interface{}{"A001"}, user)
```

On the provider side, a `context.Context` is passed as the first parameter along with the user data. Here below the code from the current samples shows how the service method is implemented, and how the user data is fetched from the attachment:

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

### How to run

#### Quick Start with makefile

*Prerequisite: docker environment*

1. **Start server**
    ```bash
    cd attachment/server
    make -f ../../build/Makefile docker-up start
    ```
2. **Run client**
    ```bash
   cd attachment/client
   make -f ../../build/Makefile run 
   ```
3. **Integration test**
   ```bash
   cd attachment/server
   make -f ../../build/Makefile integration
   ```
4. **Shutdown and cleanup**
   ```bash
   cd attachment/server
   make -f ../../build/Makefile clean docker-down
   ```

#### Run in IDE

Use *Intellij GoLand* as an example. After open dubbo-go-samples in GoLand, follow the steps below to run/debug this example:

1. **Start up zookeeper server**
   
   Open "attachment/go-server/docker/docker-compose.yaml", and click ▶︎▶︎ icon in the gutter on the left side of the editor, then "Services" tab should pop up and shows the similar message below:
   ```
   Deploying 'Compose: docker'...
   /usr/local/bin/docker-compose -f .../dubbo-go-samples/attachment/go-server/docker/docker-compose.yml up -d
   Creating network "docker_default" with the default driver
   Creating docker_zookeeper_1 ...
   'Compose: docker' has been deployed successfully.
   ```
   
2. **Start up service provider**

   Open "attachment/go-server/cmd/server.go", and click ▶︎ icon just besides "main" function in the gutter on the left side, and select "Modify Run Configuration..." from the pop-up menu. Then make sure the following configs configured correctly:
   * Working Directory: the absolute path to "attachment/go-server", for examples: */home/dubbo-go-samples/attachment/go-server*
   * Environment: CONF_PROVIDER_FILE_PATH=conf/server.yml, optionally you could also specify logging configuration with "APP_LOG_CONF_FILE=conf/log.yml"
   
   Then the sample server is ready to run. 
     
3. **Run service consumer**

   Open "attachment/go-client/cmd/client.go", and click ▶︎ icon just besides "main" function in the gutter on the left side, and select "Modify Run Configuration..." from the pop-up menu. Then make sure the following configs configured correctly:
   * Working Directory: the absolute path to "attachment/go-client", for examples: */home/dubbo-go-samples/attachment/go-client*
   * Environment: CONF_CONSUMER_FILE_PATH=conf/client.yml, optionally you could also specify logging configuration with "APP_LOG_CONF_FILE=conf/log.yml"
     
   Then run it to call the remote service, you will observe the following message output:
   ```
   [2021-02-03/16:19:30 main.main: client.go: 66] response result: &{A001 Alex Stocks 18 2020-02-04 16:19:30.422 +0800 CST}
   ```

If you need to debug either the samples or dubbo-go, you may consider switch to **Debug** instead of **Run** in GoLand. To stop, simply click ◼︎ to shutdown everything.

