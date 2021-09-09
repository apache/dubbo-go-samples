# Generic Reference Example

### 1. Introduction

Generic invocation is mainly used when the client does not have API interface or model class, all POJOs in parameters
and return values are represented by Map. Commonly used for framework integration such as: implementing a common service
testing framework, all service implementations can be invoked via
`GenericService`.

### 2. Code

```go
var ( 
    //appName is the unique identification of RPCService 
    appName         = "UserConsumer"
    referenceConfig = config.ReferenceConfig{
        InterfaceName: "org.apache.dubbo.UserProvider",
        Cluster:       "failover",
        Registry:      "demoZk",
        Protocol:      dubbo.DUBBO,
        Generic:       true,
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
    		//method name
    		"queryUser",
    		//parameter type
    		[]string{"org.apache.dubbo.User"},
    		//parameter array
    		[]hessian.Object{user},
    	},
    )
}

```
