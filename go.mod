module github.com/apache/dubbo-go-samples

require (
	dubbo.apache.org/dubbo-go/v3 v3.0.0-rc3.0.20211106163239-7480efa02126
	github.com/apache/dubbo-go-hessian2 v1.9.4-0.20210917102639-74a8ece5f3cb
	github.com/dubbogo/gost v1.11.19
	github.com/dubbogo/net v0.0.4
	github.com/dubbogo/triple v1.0.10-0.20211106162720-597a66a97296
	github.com/golang/protobuf v1.5.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pkg/errors v0.9.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/dubbogo/triple => github.com/dubbogo/triple v1.0.10-0.20211106170050-c097dda15c59

go 1.13
