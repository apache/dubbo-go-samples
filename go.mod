module github.com/apache/dubbo-go-samples

require (
	dubbo.apache.org/dubbo-go/v3 v3.0.0-20210607024529-7ef3eff3cc00
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/apache/dubbo-getty v1.4.3
	github.com/apache/dubbo-go-hessian2 v1.9.2
	github.com/bwmarrin/snowflake v0.3.0
	github.com/dubbogo/gost v1.11.9
	github.com/dubbogo/triple v1.0.0
	github.com/emicklei/go-restful/v3 v3.4.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.2.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/transaction-wg/seata-golang v0.2.0
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/willf/bitset v1.1.10 // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.26.0
)

replace (
	github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.8.0
	github.com/nacos-group/nacos-sdk-go => github.com/nacos-group/nacos-sdk-go v1.0.7-0.20210325111144-d75caca21a46
	github.com/shirou/gopsutil => github.com/shirou/gopsutil v0.0.0-20181107111621-48177ef5f880
)

go 1.13
