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

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/apache/dubbo-go-samples/task/shop/comment/api"
)

var grpcImpl = new(api.CommentClientImpl)

// export DUBBO_GO_CONFIG_PATH=../conf/dubbogo.yaml
func main() {
	config.SetConsumerService(grpcImpl)
	if err := config.Load(); err != nil {
		panic(err)
	}

	fmt.Println("start to test dubbo")
	req := &api.CommentReq{
		ItemName: "comment test",
	}
	reply, err := grpcImpl.GetComment(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply)
}
