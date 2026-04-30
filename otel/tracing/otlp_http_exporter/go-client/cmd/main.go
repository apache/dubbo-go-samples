package main

import (
	"context"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/otel/trace"

	greet "github.com/apache/dubbo-go-samples/otel/tracing/stdout/proto"

	"github.com/dubbogo/gost/log/logger"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_otel_client"),
		dubbo.WithTracing(
			trace.WithEnabled(),
			trace.WithOtlpHttpExporter(),
			trace.WithW3cPropagator(),
			trace.WithAlwaysMode(),
			trace.WithInsecure(),
			trace.WithEndpoint("127.0.0.1:4318"),
		),
	)
	if err != nil {
		panic(err)
	}

	// Triple
	cli1, err := ins.NewClient(client.WithClientURL("127.0.0.1:20000"))
	if err != nil {
		panic(err)
	}
	svc, err := greet.NewGreetService(cli1)
	if err != nil {
		panic(err)
	}
	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "triple"})
	if err != nil {
		logger.Error(err)
	} else {
		logger.Infof("Triple response: %s", resp.Greeting)
	}

	// Dubbo
	cli2, err := ins.NewClient(
		client.WithClientURL("127.0.0.1:20001"),
		client.WithClientProtocolDubbo(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		panic(err)
	}
	conn2, err := cli2.Dial("GreetProvider")
	if err != nil {
		panic(err)
	}
	var dubboResp string
	if err := conn2.CallUnary(context.Background(), []interface{}{"hello", "new", "dubbo"}, &dubboResp, "SayHello"); err != nil {
		logger.Error(err)
	} else {
		logger.Infof("Dubbo response: %s", dubboResp)
	}

	// JSONRPC
	cli3, err := ins.NewClient(
		client.WithClientURL("127.0.0.1:20002"),
		client.WithClientProtocolJsonRPC(),
		client.WithClientSerialization(constant.JSONSerialization),
	)
	if err != nil {
		panic(err)
	}
	conn3, err := cli3.Dial("GreetProvider")
	if err != nil {
		panic(err)
	}
	var jsonrpcResp string
	if err := conn3.CallUnary(context.Background(), []interface{}{"hello", "new", "jsonrpc"}, &jsonrpcResp, "SayHello"); err != nil {
		logger.Error(err)
	} else {
		logger.Infof("JSONRPC response: %s", jsonrpcResp)
	}

	if tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		_ = tp.Shutdown(context.Background())
	}
}
