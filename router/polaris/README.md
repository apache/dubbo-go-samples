# Dubbo Go & Polaris Router Example

English | [中文](README_zh.md)

## Using the service routing function

Quickly experience Polaris' service routing capabilities in dubbogo

## Polaris server installation

[Polaris Server Standalone Version Installation Documentation](https://polarismesh.cn/docs/%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%AE%89%E8%A3%85/%E5%8D%95%E6%9C%BA%E7%89%88%E5%AE%89%E8%A3%85/)

[Polaris Server Cluster Version Installation Documentation](https://polarismesh.cn/docs/%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%AE%89%E8%A3%85/%E9%9B%86%E7%BE%A4%E7%89%88%E5%AE%89%E8%A3%85/)

## How to use

[Polaris Service Routing Usage Document](https://polarismesh.cn/docs/%E5%8C%97%E6%9E%81%E6%98%9F%E6%98%AF%E4%BB%80%E4%B9%88/%E5%8A%9F%E8%83%BD%E7%89%B9%E6%80%A7/%E6%B5%81%E9%87%8F%E7%AE%A1%E7%90%86/#%E5%8A%A8%E6%80%81%E8%B7%AF%E7%94%B1)

### How to configure service routing parameters

The implementation of the PolarisMesh PriorityRouter extension point in dubbogo can automatically identify the request label information that needs to participate in service routing from the current RPC call context and request information according to the service route rules configured by the user.

![](image/dubbogo-route-rule.png)

### Running the service provider

Enter the cmd directory of server-prod, server-pre, and server-dev respectively, and execute the following commands:

```bash
# Dev server (port 20000)
cd go-server/server-dev/cmd
go run .

# Pre server (port 21000)
cd go-server/server-pre/cmd
go run .

# Prod server (port 22000)
cd go-server/server-prod/cmd
go run .
```

When you see the following log, it means that the server side started successfully:

```
INFO ... dubbo server started
```

### Running the service caller

Enter the cmd directory of go-client and execute the following commands:

```bash
cd go-client/cmd

# No uid (routes to prod)
export uid=
go run .

# uid=user-1 (routes to pre)
export uid=user-1
go run .

# uid=user-2 (routes to dev)
export uid=user-2
go run .
```

When you see the following log, it means that go-client successfully discovered go-server and made an RPC call:

```
INFO ... uid=, response: user:"[Prod] Alex Stocks"
```

### Expected output

When routing is working correctly, you will see:

```
# export uid=
uid=, response: user:"[Prod] Alex Stocks"
uid=, response: user:"[Prod] Alex Stocks"

# export uid=user-1
uid=user-1, response: user:"[Pre] Alex Stocks"
uid=user-1, response: user:"[Pre] Alex Stocks"

# export uid=user-2
uid=user-2, response: user:"[Dev] Alex Stocks"
uid=user-2, response: user:"[Dev] Alex Stocks"
```
