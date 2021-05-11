# Chain Sample

### Backend

Most of the samples in this project uses one consumer and one provider for simplification purpose. That is, there are only two nodes in the calling chain. This sample demonstrates how a calling chain which contains nodes more than two - three nodes at minimal, is configured. This sample contains three parts, which are:

frontend -> middle -> backend

1. backend: Backend services, including 'CatService', 'DogService', 'TigerService' and 'LionService'. In "backend" directory, only service providers are provided.
2. middle: Middle services, in the middle of the calling chain, provides services 'ChineseService' and 'AmericanService', and consumes the services provided by "backend" directory. In  "middle" directory, both service providers and service consumers are provided.
3. frontend: Frontend caller, consumes services provided by "middle" and output the result.   

### Call other service in the current service

**Code**

```golang
type DogService struct {
	GetID   func() (int, error)
	GetName func() (string, error)
	Yell    func() (string, error)
}

func (d *DogService) Reference() string {
	return "DogService"
}

type TigerService struct {
	GetID   func() (int, error)
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

**Configuration**

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

Pls. refer to [HOWTO.md](../HOWTO.md) under the root directory to run this sample.
