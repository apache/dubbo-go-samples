module github.com/apache/dubbo-go-samples

require (
	cloud.google.com/go v0.50.0 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.3.0
	github.com/Microsoft/go-winio v0.4.15-0.20190919025122-fc70bd9a86b5 // indirect
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/apache/dubbo-getty v1.4.3
	github.com/apache/dubbo-go v1.5.7-rc1
	github.com/apache/dubbo-go-hessian2 v1.9.2
	github.com/bwmarrin/snowflake v0.3.0
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/dubbogo/gost v1.11.2
	github.com/emicklei/go-restful/v3 v3.4.0
	github.com/golang/protobuf v1.4.3
	github.com/gophercloud/gophercloud v0.3.0 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/linode/linodego v0.10.0 // indirect
	github.com/miekg/dns v1.1.27 // indirect
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/opentrx/mysql v1.0.0-pre // indirect
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/transaction-wg/seata-golang v0.2.1-alpha
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.opencensus.io v0.23.0
	google.golang.org/grpc v1.33.2
	google.golang.org/grpc/examples v0.0.0-20210322221411-d26af8e39165 // indirect
)

replace (
	github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.8.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

go 1.13
