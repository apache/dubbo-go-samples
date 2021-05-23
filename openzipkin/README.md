# Zipkin in Dubbo-go Example

## Backend

Zipkin is a distributed tracing system. It helps gather timing data needed to troubleshoot latency problems in service architectures. Features include both the collection and lookup of this data.

## Introduction



```markdown
.
├── README.md
├── README_zh.md
├── docker-compose.yml
├── go-client
├── go-server-a
├── go-server-b
└── prometheus
```

- go-client ：The Service Consumer
- go-server-a ：The Service Provider A
- go-server-b ：The Service Provider B

## Code

register Zipkin. Reporter，Endpoint，Tracer，default sample `AlwaysSample` 

```go
func registerZipkin() {
	// set up a span reporter
	reporter := zipkinhttp.NewReporter("http://localhost:9411/api/v2/spans")

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint("go-server-a", "localhost:80")
	if err != nil {
		gxlog.CError("unable to create local endpoint: %+v\n", err)
	}

    // set sampler , default AlwaysSample
    // sampler := zipkin.NewModuloSampler(1)

	// initialize our tracer
	// nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint), zipkin.WithSampler(sampler))
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		gxlog.CError("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)
}
```

## Config

Provider config filter：

```yaml

services:
  ...
filter: "tracing"

```

## Filter

Dubbo-go support `opentrace filter` 实现，基于简单配置即可

```go
func (tf *tracingFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	var (
		spanCtx context.Context
		span    opentracing.Span
	)
	operationName := invoker.GetUrl().ServiceKey() + "#" + invocation.MethodName()

	wiredCtx := ctx.Value(constant.TRACING_REMOTE_SPAN_CTX)
	preSpan := opentracing.SpanFromContext(ctx)

	if preSpan != nil {
		// it means that someone already create a span to trace, so we use the span to be the parent span
		span = opentracing.StartSpan(operationName, opentracing.ChildOf(preSpan.Context()))
		spanCtx = opentracing.ContextWithSpan(ctx, span)

	} else if wiredCtx != nil {

		// it means that there has a remote span, usually from client side. so we use this as the parent
		span = opentracing.StartSpan(operationName, opentracing.ChildOf(wiredCtx.(opentracing.SpanContext)))
		spanCtx = opentracing.ContextWithSpan(ctx, span)
	} else {
		// it means that there is not any span, so we create a span as the root span.
		span, spanCtx = opentracing.StartSpanFromContext(ctx, operationName)
	}

	defer func() {
		span.Finish()
	}()

	result := invoker.Invoke(spanCtx, invocation)
	span.SetTag(successKey, result.Error() == nil)
	if result.Error() != nil {
		span.LogFields(log.String(errorKey, result.Error().Error()))
	}
	return result
}
```

## Install Zipkin

First Way:

Follow [Zipkin's quick start](https://zipkin.io/pages/quickstart.html) to install zipkin.

```bash
curl -sSL https://zipkin.io/quickstart.sh | bash -s
```

Zipkin supports various backend storages including Cassandra, ElasticSearch and MySQL. Here we use the simplest storage - in-memory for demo purpose.

```bash
java -jar zipkin.jar
```

Once the process starts, you can verify zipkin server works by access http://localhost:9411

Or Use Docker:

See [/dubbo-go-sample/zipkin/docker/docker-compose.yml](docker-compose.yml)

```dockerfile
services:
  # The zipkin process services the UI, and also exposes a POST endpoint that
  # instrumentation can send trace data to.
  zipkin:
    image: ghcr.io/openzipkin/zipkin-slim:${TAG:-latest}
    container_name: zipkin
    # Environment settings are defined here https://github.com/openzipkin/zipkin/blob/master/zipkin-server/README.md#environment-variables
    environment:
      - STORAGE_TYPE=mem
      # Point the zipkin at the storage backend
      - MYSQL_HOST=mysql
      # Uncomment to enable self-tracing
      # - SELF_TRACING_ENABLED=true
      # Uncomment to increase heap size
      # - JAVA_OPTS=-Xms128m -Xmx128m -XX:+ExitOnOutOfMemoryError
    ports:
      # Port used for the Zipkin UI and HTTP Api
      - 9411:9411
    # Uncomment to enable debug logging
    # command: --logging.level.zipkin2=DEBUG
```

```bash
docker-compose -f docker/docker-compose.yml up -d zipkin
```

### How To Run

Refer to  [HOWTO.md](../HOWTO_zh.md) under the root directory to run this sample.


