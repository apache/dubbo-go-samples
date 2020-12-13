module github.com/apache/dubbo-go-samples

require (
	github.com/apache/dubbo-go v1.5.2-0.20201213072836-1f708d7cb6c5
	github.com/apache/dubbo-go-hessian2 v1.7.0
	github.com/bwmarrin/snowflake v0.3.0
	github.com/dk-lockdown/seata-golang v0.1.0-fix
	github.com/dubbogo/gost v1.9.5
	github.com/emicklei/go-restful/v3 v3.4.0
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0
	github.com/stretchr/testify v1.6.1
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	google.golang.org/grpc v1.26.0
)

replace github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.8.0

go 1.13
