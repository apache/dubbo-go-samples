module github.com/apache/dubbo-go-samples

require (
	cloud.google.com/go v0.50.0 // indirect
	github.com/Microsoft/go-winio v0.4.15-0.20190919025122-fc70bd9a86b5 // indirect
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/apache/dubbo-go v1.5.6-rc1.0.20210225142527-f1f81ebd8292
	github.com/apache/dubbo-go-hessian2 v1.9.1
	github.com/bwmarrin/snowflake v0.3.0
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/dubbogo/gost v1.11.2
	github.com/emicklei/go-restful/v3 v3.4.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/gophercloud/gophercloud v0.3.0 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/linode/linodego v0.10.0 // indirect
	github.com/miekg/dns v1.1.27 // indirect
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/transaction-wg/seata-golang v0.1.8
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.etcd.io/bbolt v1.3.4 // indirect
	google.golang.org/grpc v1.36.0
	google.golang.org/grpc/examples v0.0.0-20210514222045-a12250e98f97 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace (
	github.com/apache/dubbo-go => github.com/apache/dubbo-go v1.5.7-rc1-tmp
	github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.8.0
)

go 1.13
