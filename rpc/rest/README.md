# Dubbo-Go REST sample

This sample validates the Dubbo-Go REST protocol with direct URL, interface-level registry, and application-level service discovery. It uses an explicit REST mapping for path, query, header, and body arguments.

## Run

Start the provider with the default Nacos application-level service discovery mode:

```bash
go run ./rpc/rest/go-server/cmd
```

Run the Dubbo-Go REST consumer in another terminal:

```bash
go run ./rpc/rest/go-client/cmd
```

Expected client output:

```text
REST response: userID=101 name=dubbo-go traceID=trace-rest-basic message=body-from-dubbo-rest-client greeting="hello dubbo-go, userID=101, traceID=trace-rest-basic, message=body-from-dubbo-rest-client"
```

The same provider is also a normal HTTP endpoint:

```bash
curl -s \
  -X POST 'http://127.0.0.1:20080/api/v1/users/202/greeting?name=curl' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'X-Trace-ID: trace-curl' \
  -d '{"message":"body-from-curl"}'
```

## What this proves

The provider URL only supplies the network target, for example `rest://127.0.0.1:20080/org.apache.dubbo.samples.rest.GreetingService`.

The actual REST call shape comes from the REST method config in `api/rest_config.go`:

- method: `POST`
- path: `/api/v1/users/{userID}/greeting`
- path parameter: argument `0` -> `userID`
- query parameter: argument `1` -> `name`
- header: argument `2` -> `X-Trace-ID`
- body: argument `3`

So this sample intentionally separates "provider address" from "REST HTTP mapping".

## Registry service discovery

The provider and consumer support these flags:

- `-registry=direct|zookeeper|nacos`
- `-registry-type=interface|service|all`

The default is `-registry=nacos -registry-type=service`.

`interface` means the registry stores the callable `rest://...` provider URL directly. `service` means the registry stores the application instance; the consumer finds the application by service mapping, fetches metadata from the provider metadata service, then reconstructs the `rest://...` provider URL from the discovered instance and service metadata. `all` registers both forms.

Run direct URL mode:

```bash
go run ./rpc/rest/go-server/cmd -registry=direct
go run ./rpc/rest/go-client/cmd -registry=direct
```

Run ZooKeeper interface-level registration:

```bash
go run ./rpc/rest/go-server/cmd -registry=zookeeper -registry-type=interface
go run ./rpc/rest/go-client/cmd -registry=zookeeper -registry-type=interface
```

All modes should print the same `REST response: ...` line from the consumer.
