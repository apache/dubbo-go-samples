module github.com/apache/dubbo-go-samples

require (
	dubbo.apache.org/dubbo-go/v3 v3.0.0-rc2.0.20210905041349-a1e4548f1e68
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/apache/dubbo-getty v1.4.5
	github.com/apache/dubbo-go-hessian2 v1.9.2
	github.com/bwmarrin/snowflake v0.3.0
	github.com/dubbogo/gost v1.11.16
	github.com/dubbogo/net v0.0.4
	github.com/dubbogo/triple v1.0.6-0.20210904050749-5721796f3fd6
	github.com/emicklei/go-restful/v3 v3.5.2
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/transaction-wg/seata-golang v0.2.0
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.27.1
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a // indirect
)

replace github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.8.0

go 1.13
