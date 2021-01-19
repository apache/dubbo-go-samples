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
	"fmt"
	pb "github.com/apache/dubbo-go-samples/general/dubbo3/protobuf/grpc"
	"log"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:19998"
	defaultName = "jifeng"
)

func main() {
	// Set up a connection to the client.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDubbo3GreeterClient(conn)

	// Contact the client and print out its response.
	//name := defaultName
	//if len(os.Args) > 1 {
	//	name = os.Args[1]
	//}
	BigDataReq := pb.BigData{
		WantSize: 271828,
		Data: make([]byte, 314159),
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			// Test Big Unary
			rsp, err := c.BigUnaryTest(context.Background(), &BigDataReq)
			if err != nil{
				panic(err)
			}
			fmt.Println("rsp len = ", len(rsp.Data))
			//r, err := c.Dubbo3SayHello2(context.Background(), &pb.Dubbo3HelloRequest{Myname: name})
			//if err != nil {
			//	fmt.Print("could not greet: %v", err)
			//}
			//log.Printf("####### get client %+v", r)
			//wg.Done()
		}()
	}
	wg.Wait()
}
