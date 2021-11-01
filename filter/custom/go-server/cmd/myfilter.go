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
	extension.SetFilter("myServerFilter", NewMyServerFilter)
}

func NewMyServerFilter() filter.Filter {
	return &MyServerFilter{}
}

type MyServerFilter struct {
}

func (f *MyServerFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	fmt.Println("MyServerFilter Invoke is called, method Name = ", invocation.MethodName())
	return invoker.Invoke(ctx, invocation)
}
func (f *MyServerFilter) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, protocol protocol.Invocation) protocol.Result {
	fmt.Println("MyServerFilter OnResponse is called")
	return result
}
