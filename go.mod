module github.com/apache/dubbo-go-samples

require (
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/apache/dubbo-go v1.5.6-0.20210204121013-9c3bd0e9c92d
	github.com/apache/dubbo-go-hessian2 v1.8.2
	github.com/bwmarrin/snowflake v0.3.0
	github.com/transaction-wg/seata-golang v0.2.0
	github.com/dubbogo/gost v1.11.0
	github.com/emicklei/go-restful/v3 v3.4.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/labstack/echo/v4 v4.1.15 // indirect
	github.com/micro/go-micro/v2 v2.9.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	google.golang.org/grpc v1.26.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1 // indirect
)

replace github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.8.0

go 1.13
