## Hello World

Hello World Demo

### Configuration

registy config

```yaml
# registry config
registries:
  "demoZk":
    protocol: "zookeeper"
    timeout: "3s"
    address: "127.0.0.1:2181"

```

provider config

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

consumer config

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

### Code

provider service define

```
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

consumer reference define

```
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


### Action

1. Run zookeeper with docker

```shell
cd helloworld/go-server
make -f ../../build/Makefile docker-up
```

2. Run provider application

```shell
cd helloworld/go-server
make -f ../../build/Makefile start
```

3. Run consumer application

```shell
cd helloworld/go-client
make -f ../../build/Makefile run
```