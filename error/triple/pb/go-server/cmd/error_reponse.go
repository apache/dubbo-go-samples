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
)

import (
	"github.com/dubbogo/gost/log/logger"

	"github.com/dubbogo/grpc-go/codes"
	"github.com/dubbogo/grpc-go/status"
)

import (
	triplepb "github.com/apache/dubbo-go-samples/api"
)

type ErrorResponseProvider struct {
	triplepb.UnimplementedGreeterServer
}

func (s *ErrorResponseProvider) SayHello(ctx context.Context, in *triplepb.HelloRequest) (*triplepb.User, error) {
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n" + in.Name)

	/* GRPC/Triple wrapped error, client would get:
	 error details = [type.googleapis.com/google.rpc.DebugInfo]:{stack_entries:"
	 main.(*ErrorResponseProvider).SayHello
	       xxx/dubbo-go-samples/error/triple/pb/go-server/cmd/error_reponse.go:48
	...
	 error code = Code(1234)
	 error message = user defined error
	*/
	return &triplepb.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, status.Error(codes.Code(1234), "user defined error")

	/* normal error with stack, client would get:
	error details = [type.googleapis.com/google.rpc.DebugInfo]:{stack_entries:"userDefinedError
	main.(*ErrorResponseProvider).SayHello
	       xxx/dubbo-go-samples/error/triple/pb/go-server/cmd/error_reponse.go:55
	error code = Unknown
	error message = userDefinedError
	*/
	//return &triplepb.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, errors.New("userDefinedError")
}
