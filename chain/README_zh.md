# 链路调用示例

### 背景

为了简化起见，绝大部分的例子都是展示了一个 consumer 调用一个 provider 的场景。也就是说，这个调用链路上只有两个节点。本例展示了大于两个节点的最小链路单元 —— 三节点，dubbo-go 是如何配置和工作的。在这个例子中，由三部分组成：

frontend -> middle -> backend

1. backend: 后端服务，包含 CatService、DogService、TigerService 和 LionService。"backend" 目录只包含了服务的提供者。
2. middle: 中间服务，在链路的中间，负责向前端（frontend）暴露服务 ChineseService 和 AmericanService，并调用后端（backend）服务。"middle" 目录既包含了服务的提供者，也包含了服务的调用者。
3. frontend: 前端调用者，调用 middle 提供的服务并输出。

### 在提供的服务中调用其他服务的示例

**代码**

```golang
type DogService struct {
	GetId   func() (int, error)
	GetName func() (string, error)
	Yell    func() (string, error)
}

func (d *DogService) Reference() string {
	return "DogService"
}

type TigerService struct {
	GetId   func() (int, error)
	GetName func() (string, error)
	Yell    func() (string, error)
}

func (t *TigerService) Reference() string {
	return "TigerService"
}

func init() {
	dog := new(DogService)
	config.SetConsumerService(dog)
	tiger := new(TigerService)
	config.SetConsumerService(tiger)

	config.SetProviderService(&ChineseService{
		dog:   dog,
		tiger: tiger,
	})
}
```

**配置**

```yaml
# reference config
references:
  "CatService":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.demo.CatService"
  "DogService":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.demo.DogService"
    
# service config
services:
  "ChineseService":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.demo.ChineseService"
  "AmericanService":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.demo.AmericanService"
```


请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。


