## Hello World 实例

### 配置

注册中心配置

```yaml
# registry config
registries:
  "demoZk":
    protocol: "zookeeper"
    timeout: "3s"
    address: "127.0.0.1:2181"

```

服务提供者配置

```yaml
# service config
services:
  # Reference ID
  "UserProvider":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    cluster: "failover"
    methods:
      - name: "GetUser"
        retries: 1
```

服务消费者配置

```yaml
# reference config
references:
  # Reference ID
  "UserProvider":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    cluster: "failover"
    methods:
      - name: "GetUser"
        retries: 3
```

### 代码示例

生产者示例

```go
// init 
func init() {
	config.SetProviderService(new(UserProvider))
	// ------for hessian2------
	hessian.RegisterPOJO(&User{})
}

// define dto
type User struct {
	Id   string
	Name string
	Age  int32
	Time time.Time
}

// implement POJO interface for hessian2
func (u User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

// service define
type UserProvider struct {
}

// interface define
func (u *UserProvider) GetUser(ctx context.Context, req []interface{}) (*User, error) {
	//biz code...
}

// implement RPCService interface
func (u *UserProvider) Reference() string {
	return "UserProvider"
}
```

消费者示例

```go
var userProvider = new(pkg.UserProvider)

// init 
func init() {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&pkg.User{})
}

// define dto
type User struct {
	Id   string
	Name string
	Age  int32
	Time time.Time
}

// implement POJO interface for hessian2
func (u User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

// service define
type UserProvider struct {
    GetUser func(ctx context.Context, req []interface{}, rsp *User) error
}

// implement RPCService interface
func (u *UserProvider) Reference() string {
	return "UserProvider"
}

func main() {
    //dubbogo init
    config.Load()
    time.Sleep(3 * time.Second)
    
    user := &pkg.User{}
    err := userProvider.GetUser(context.TODO(), []interface{}{"A001"}, user)
    if err != nil {
        //...
    }
    gxlog.CInfo("response result: %v\n", user)
}
```