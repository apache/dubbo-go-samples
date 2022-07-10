module github.com/apache/dubbo-go-samples

require (
	dubbo.apache.org/dubbo-go/v3 v3.0.3-0.20220708050457-1abbc2e77c76
	github.com/apache/dubbo-go-hessian2 v1.11.0
	github.com/dubbogo/gost v1.12.5
	github.com/dubbogo/grpc-go v1.42.9
	github.com/dubbogo/triple v1.1.8
	github.com/golang/protobuf v1.5.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/jaeger v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
)

replace dubbo.apache.org/dubbo-go/v3 => /Users/phil/Pi/cygnusspace/dubbo-go

go 1.15
