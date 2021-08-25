## 通过配置API生成实例化子模块

### 注册中心

例子：
```go
// 通过配置 API 生成注册中心配置，此处为默认zk配置
registryConfig := config.NewRegistryConfigWithProtocolDefaultPort("zookeeper")

// 调用注册中心配置API，获取注册中心实例 reg
reg, err := registryConfig.GetInstance(common.PROVIDER)
if err != nil {
    panic(err)
}

// 生成一个 provider URL
ivkURL, err := common.NewURL("mock://localhost:8080",
    common.WithPath("com.alibaba.dubbogo.HelloService"),
    common.WithParamsValue(constant.ROLE_KEY, strconv.Itoa(common.PROVIDER)),
    common.WithMethods([]string{"GetUser", "SayHello"}),
)
if err != nil {
    panic(err)
}
// 使用新生成的注册中心实例注册
if err := reg.Register(ivkURL); err != nil {
    panic(err)
}
time.Sleep(time.Second * 30)
// 反注册
if err := reg.UnRegister(ivkURL); err != nil {
    panic(err)
}

```