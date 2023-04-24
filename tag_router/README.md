# dubbo-go tag router example

## 1 Introduction

Tag router is an important part of Dubbo traffic management. Tag router achieves traffic isolation by dividing an instance of a service into different groups, constraining traffic with specific tags to only flow within specified groups. Different groups serve different traffic scenarios, and can serve as the basis for scenarios such as blue green publishing and gray publishing.

## 2 Usage
### 2.1 provider

There are two ways to label the service provider on the provider side, static labeling and dynamic labeling.

Static labeling is the process of labeling service providers in the configuration file before service startup. Every time a label is changed, the service provider needs to be restarted.

As shown below, the service has been labeled as tag1. **Note: The object being labeled is at the application level**
```yaml
dubbo:
  application:
    tag: tag1
  registries:
    demoZK:
      protocol: zookeeper
      address: 127.0.0.1:2181
  protocols:
    triple:
      name: tri
      port: 20000
  provider:
    services:
      GreeterProvider:
        interface: com.apache.dubbo.sample.basic.IGreeter # must be compatible with grpc or dubbo-java
```
Dynamic labeling is more flexible than static marking, allowing for the replacement of service labels during the service processing. But users need to operate through dubbo admin.

Enter the tag routing module under dubbo admin service governance. Clicking Create will pop up the Create tag Routing window. Relevant rules can be configured in the rule content.
```yaml
force: false
enabled: true
runtime: false
tags:
  - name: tag1
    addresses: [192.168.0.1:20881]
  - name: tag2
    addresses: [192.168.0.2:20882]
```
After the configuration is completed, click Save to dynamically modify the label of the service provider.

### 2.2 consumer

On the consumer side, you can choose which service provider to be used by tag, defined in the code. The following are the services provided using tag1.
```go
// set tag
ctx := context.Background()
atm := map[string]string{
    "dubbo.tag":       "tag1",
    "dubbo.force.tag": "true",
}
ctx = context.WithValue(ctx, constant.AttachmentKey, atm)
reply, err := grpcGreeterImpl.SayHello(ctx, req)
if err != nil {
    logger.Error(err)
}
logger.Infof("client response result: %v\n", reply)
```

## 3 Deploy

Users can deploy this demo to try using tag router, and Dubbo go needs to be upgraded to the latest version.
1. Start the zookeeper and dubbo admin;
2. Set the configuration file environment variables of the provider and start the provider;
3. Set the configuration file environment variables of the consumer and start the consumer;
4. Reboot after modifying the provider's tag through configuration locally;
5. Attempt to modify the provider's tag through the dubbo admin.