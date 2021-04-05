module github.com/apache/dubbo-go-samples

require (
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/apache/dubbo-go v1.5.6-rc2.0.20210405062714-19c3cc3c36b1
	github.com/apache/dubbo-go-hessian2 v1.9.1
	github.com/bwmarrin/snowflake v0.3.0
	github.com/dubbogo/gost v1.11.3
	github.com/dubbogo/v3router v0.1.0
	github.com/emicklei/go-restful/v3 v3.4.0
	github.com/golang/protobuf v1.5.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/transaction-wg/seata-golang v0.2.0
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	google.golang.org/grpc v1.33.1
	google.golang.org/grpc/examples v0.0.0-20210322221411-d26af8e39165 // indirect
)

replace github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.8.0

replace github.com/nacos-group/nacos-sdk-go => github.com/nacos-group/nacos-sdk-go v1.0.7-0.20210325111144-d75caca21a46

go 1.13
