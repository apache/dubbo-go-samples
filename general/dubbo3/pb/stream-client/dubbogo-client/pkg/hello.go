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

package pkg

import (
	"context"
)

import (
	"github.com/dubbogo/triple/pkg/triple"
)

import (
	dubbo3pb "github.com/apache/dubbo-go-samples/general/dubbo3/pb/protobuf/dubbo3"
)

type GreeterProvider struct {
	SayHelloStream func(ctx context.Context) (dubbo3pb.Greeter_SayHelloStreamClient, error)
}

func (u *GreeterProvider) Reference() string {
	return "GreeterProvider"
}

func (u *GreeterProvider) GetDubboStub(cc *triple.TripleConn) dubbo3pb.GreeterClient {
	return dubbo3pb.NewGreeterDubbo3Client(cc)
}
