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

package integration

import (
	"os"
	"testing"

	"github.com/dubbogo/gost/log/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/apache/dubbo-go-samples/rpc/grpc/proto"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

var greetService pb.GreetServiceClient

func TestMain(m *testing.M) {
	dubboConn, err := grpc.Dial("127.0.0.1:20000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer dubboConn.Close()
	greetService = pb.NewGreetServiceClient(dubboConn)
	os.Exit(m.Run())
}
