# 泛化调用实例

### 1. 介绍

泛化接口调用方式主要用于客户端没有 API 接口及模型类元的情况，参数及返回值中的所有 POJO 均用 Map 表示，通常用于框架集成，比如：实现一个通用的服务测试框架，可通过`GenericService`调用所有服务实现。

### 2. 代码示例

```go
var ( 
    //appName是RPCService的唯一标识 
    appName         = "UserConsumer"
    referenceConfig = config.ReferenceConfig{
        InterfaceName: "org.apache.dubbo.UserProvider",
        Cluster:       "failover",
        Registry:      "demoZk",
        Protocol:      dubbo.DUBBO,
        Generic:       "true",
    }
)

func init() {
    referenceConfig.GenericLoad(appName)
    time.Sleep(3 * time.Second)
}

func main() {
    resp, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
    	context.TODO(),
    	[]interface{}{
    		//方法名
    		"queryUser",
    		//参数类型数组
    		[]string{"org.apache.dubbo.User"},
    		//参数数组
    		// user = &User{
    		//  ID: "3213",
    		//  Name: "panty",
    		//  Age: 25,
    		//  Time: time.Now()
    		// } 
    		[]interface{}{
    			map[string]interface{}{
                    "iD": "3213",
                    "name": "panty",
                    "age": 25,
                    "time": time.Now(),
    			},
    		},
    	},
    )
}

```