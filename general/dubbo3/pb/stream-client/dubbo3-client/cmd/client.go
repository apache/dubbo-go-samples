/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	protobuf "github.com/apache/dubbo-go-samples/general/dubbo3/pb/protobuf/dubbo3"
	"go.uber.org/atomic"
	"sync"
	"time"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo3"
	_ "github.com/apache/dubbo-go/protocol/grpc"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
	_ "github.com/apache/dubbo-go/registry/protocol"
	"github.com/apache/dubbo-go-samples/general/dubbo3/pb/stream-client/dubbo3-client/pkg"
	dubbo3 "github.com/apache/dubbo-go-samples/general/dubbo3/pb/protobuf/dubbo3"
)

var grpcGreeterImpl = new(pkg.GrpcGreeterImpl)

func init() {
	config.SetConsumerService(grpcGreeterImpl)
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	config.Load()
	time.Sleep(time.Second * 3)
	//testStreamClient()
	//testMultiThreadStreamClient()
	wg := sync.WaitGroup{}
	for i := 0; i < 200;  i++{
		wg.Add(1)
		go func() {
			testBigData()
			wg.Done()
		}()
	}
	wg.Wait()
}

var firstSend atomic.Int32
var secondSend atomic.Int32
var thirdSend atomic.Int32
var firstRecv atomic.Int32
var secondRecv atomic.Int32


func testBigData(){
	ctx := context.Background()
	//ctx = context.WithValue(ctx,"tri-req-id","id value" )
	dataSend := make([]byte, 300000)
	copy( dataSend[:5] , []byte("hello"))
	copy(dataSend[len(dataSend) - 5:] , []byte("world"))
	req := dubbo3.BigData{
		WantSize: 300000,
		Data: dataSend,
	}

	r, err := grpcGreeterImpl.BigStreamTest(ctx)
	if err != nil {
		panic(err)
	}
	//start := time.Now().Nanosecond()
	for i := 0; i < 3; i++{
		if err := r.Send(&req); err != nil {
			fmt.Println("say hello err:", err)
		}
		//if i == 0{
		//	fmt.Println("firstSend = ", firstSend.Inc())
		//}
		//if i == 1{
		//	fmt.Println("secondSend = ", secondSend.Inc())
		//}
		//if i == 2{
		//	thirdSend.Inc()
		//	fmt.Println("thirdSend = ", thirdSend.Inc())
		//}
		req.WantSize++
	}

	rsp := &dubbo3.BigData{}
	if err := r.RecvMsg(rsp); err != nil {
		fmt.Println("err = ", err)
	}
	//fmt.Println("firstRecv = ", firstRecv.Inc())
	fmt.Printf("firstSend Got rsp = %+v\n", len(rsp.Data))
	rsp = &dubbo3.BigData{}
	if err := r.RecvMsg(rsp); err != nil {
		fmt.Println("err = ", err)
	}
	//fmt.Println("secondRecv = ", secondRecv.Inc())
	fmt.Printf("sendSend Got rsp = %+v\n", len(rsp.Data))
	//fmt.Println("time cost = ", time.Now().Nanosecond() - start)
}

func testMultiThreadStreamClient() {
	ctx := context.Background()
	//ctx = context.WithValue(ctx,"tri-req-id","id value" )

	wg := sync.WaitGroup{}
	clientList := make([]protobuf.Dubbo3Greeter_Dubbo3SayHelloClient, 0,1000 )
	for i := 0; i < 1; i++ {
		r, err := grpcGreeterImpl.Dubbo3SayHello(ctx)
		if err != nil {
			panic(err)
		}
		clientList = append(clientList, r)
		if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng1"}); err != nil {
			fmt.Println("say hello err:", err)
		}
		if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng2"}); err != nil {
			fmt.Println("say hello err:", err)
		}
		if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng3"}); err != nil {
			fmt.Println("say hello err:", err)
		}

		rsp := &dubbo3.Dubbo3HelloReply{}
		if err := r.RecvMsg(rsp); err != nil {
			fmt.Println("err = ", err)
		}
		fmt.Printf("firstSend Got rsp = %+v\n", rsp)
		rsp = &dubbo3.Dubbo3HelloReply{}
		if err := r.RecvMsg(rsp); err != nil {
			fmt.Println("err = ", err)
		}
		if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng1"}); err != nil {
			fmt.Println("say hello err:", err)
		}
		if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng2"}); err != nil {
			fmt.Println("say hello err:", err)
		}
		if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng3"}); err != nil {
			fmt.Println("say hello err:", err)
		}

		rsp = &dubbo3.Dubbo3HelloReply{}
		if err := r.RecvMsg(rsp); err != nil {
			fmt.Println("err = ", err)
		}
		fmt.Printf("firstSend Got rsp = %+v\n", rsp)
		rsp = &dubbo3.Dubbo3HelloReply{}
		if err := r.RecvMsg(rsp); err != nil {
			fmt.Println("err = ", err)
		}
		if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng1"}); err != nil {
			fmt.Println("say hello err:", err)
		}
		if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng2"}); err != nil {
			fmt.Println("say hello err:", err)
		}
		if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng3"}); err != nil {
			fmt.Println("say hello err:", err)
		}

		rsp = &dubbo3.Dubbo3HelloReply{}
		if err := r.RecvMsg(rsp); err != nil {
			fmt.Println("err = ", err)
		}
		fmt.Printf("firstSend Got rsp = %+v\n", rsp)
		rsp = &dubbo3.Dubbo3HelloReply{}
		if err := r.RecvMsg(rsp); err != nil {
			fmt.Println("err = ", err)
		}
		fmt.Printf("firstSend Got rsp = %+v\n", rsp)
	}
	wg.Wait()
	time.Sleep(time.Second*3)
}

func testStreamClient() {
	ctx := context.Background()
	//ctx = context.WithValue(ctx,"tri-req-id","id value" )
	r, err := grpcGreeterImpl.Dubbo3SayHello(ctx)
	if err != nil {
		panic(err)
	}
	if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng"}); err != nil {
		fmt.Println("say hello err:", err)
	}
	if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng"}); err != nil {
		fmt.Println("say hello err:", err)
	}
	if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng"}); err != nil {
		fmt.Println("say hello err:", err)
	}

	rsp := &dubbo3.Dubbo3HelloReply{}
	if err := r.RecvMsg(rsp); err != nil {
		fmt.Println("err = ", err)
	}
	fmt.Printf("firstSend Got rsp = %+v\n", rsp)
	rsp = &dubbo3.Dubbo3HelloReply{}
	if err := r.RecvMsg(rsp); err != nil {
		fmt.Println("err = ", err)
	}
	fmt.Printf("firstSend Got rsp = %+v\n", rsp)

	r, err = grpcGreeterImpl.Dubbo3SayHello(ctx)
	if err != nil {
		panic(err)
	}
	if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng"}); err != nil {
		fmt.Println("say hello err:", err)
	}
	if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng"}); err != nil {
		fmt.Println("say hello err:", err)
	}
	if err := r.Send(&dubbo3.Dubbo3HelloRequest{Myname: "jifeng"}); err != nil {
		fmt.Println("say hello err:", err)
	}

	rsp = &dubbo3.Dubbo3HelloReply{}
	if err := r.RecvMsg(rsp); err != nil {
		fmt.Println("err = ", err)
	}
	fmt.Printf("firstSend Got rsp = %+v\n", rsp)
	rsp = &dubbo3.Dubbo3HelloReply{}
	if err := r.RecvMsg(rsp); err != nil {
		fmt.Println("err = ", err)
	}
	fmt.Printf("firstSend Got rsp = %+v\n", rsp)
}
