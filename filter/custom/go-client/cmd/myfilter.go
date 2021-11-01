package main

import (
	"context"
	"fmt"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
)

func init() {
	extension.SetFilter("myClientFilter", NewMyClientFilter)
}

func NewMyClientFilter() filter.Filter {
	return &MyClientFilter{}
}

type MyClientFilter struct {
}

func (f *MyClientFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	fmt.Println("MyClientFilter Invoke is called, method Name = ", invocation.MethodName())
	return invoker.Invoke(ctx, invocation)
}
func (f *MyClientFilter) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, protocol protocol.Invocation) protocol.Result {
	fmt.Println("MyClientFilter OnResponse is called")
	return result
}
