package main

import (
	"fmt"
	"github.com/apache/dubbo-go-samples/chain/frontend/pkg"
	"github.com/apache/dubbo-go/config"
	"time"

	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/metadata/service/inmemory"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

func main() {
	var chinese = new(pkg.ChineseService)
	var american = new(pkg.AmericanService)
	config.SetConsumerService(chinese)
	config.SetConsumerService(american)
	config.Load()
	time.Sleep(3 * time.Second)
	have, _ := chinese.Have()
	fmt.Printf("chinese.Have(): %s\n", have)
	hear, _ := chinese.Hear()
	fmt.Printf("chinese.Hear(): %s\n", hear)
	have, _ = american.Have()
	fmt.Printf("american.Have(): %s\n", have)
	hear, _ = american.Hear()
	fmt.Printf("american.Hear(): %s\n", hear)
}
