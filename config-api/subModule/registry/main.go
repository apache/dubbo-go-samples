package main

import (
	"strconv"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

func main() {
	registryConfig := config.NewRegistryConfigWithProtocolDefaultPort("zookeeper")
	reg, err := registryConfig.GetInstance(common.PROVIDER)
	if err != nil {
		panic(err)
	}

	ivkURL, err := common.NewURL("mock://localhost:8080",
		common.WithPath("com.alibaba.dubbogo.HelloService"),
		common.WithParamsValue(constant.ROLE_KEY, strconv.Itoa(common.PROVIDER)),
		common.WithMethods([]string{"GetUser", "SayHello"}),
	)
	if err != nil {
		panic(err)
	}
	if err := reg.Register(ivkURL); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 30)
	if err := reg.UnRegister(ivkURL); err != nil {
		panic(err)
	}
}
