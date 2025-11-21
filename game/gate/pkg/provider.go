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
	"github.com/dubbogo/gost/log/logger"

	"google.golang.org/protobuf/types/known/structpb"
)

import (
	gateProto "github.com/apache/dubbo-go-samples/game/proto/gate"
)

type GateServiceHandler struct{}

func (p *GateServiceHandler) Send(ctx context.Context, req *gateProto.SendRequest) (*gateProto.Result, error) {
	logger.Infof("football: to=%s, message=%s", req.Uid, req.Data)
	data, _ := structpb.NewStruct(map[string]interface{}{
		"to":      req.Uid,
		"message": req.Data,
	})
	return &gateProto.Result{Code: 0, Data: data}, nil
}
