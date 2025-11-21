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
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	gameProto "github.com/apache/dubbo-go-samples/game/proto/game"
)

var GameFootball gameProto.GameService

func InitGameClient() (gameProto.GameService, error) {
	cli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
	)
	if err != nil {
		return nil, err
	}

	svc, err := gameProto.NewGameService(cli)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func SetGameClient(cli gameProto.GameService) {
	GameFootball = cli
	logger.Info("game client initialized")
}

// Helper functions for HTTP handlers
type Result struct {
	Code int32                  `json:"code"`
	Msg  string                 `json:"msg,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

func Login(ctx context.Context, data string) (*Result, error) {
	resp, err := GameFootball.Login(ctx, &gameProto.LoginRequest{Uid: data})
	if err != nil {
		return nil, err
	}
	return convertResult(resp), nil
}

func Score(ctx context.Context, uid, score string) (*Result, error) {
	resp, err := GameFootball.Score(ctx, &gameProto.ScoreRequest{Uid: uid, Score: score})
	if err != nil {
		return nil, err
	}
	return convertResult(resp), nil
}

func Rank(ctx context.Context, uid string) (*Result, error) {
	resp, err := GameFootball.Rank(ctx, &gameProto.RankRequest{Uid: uid})
	if err != nil {
		return nil, err
	}
	return convertResult(resp), nil
}

func convertResult(resp *gameProto.Result) *Result {
	result := &Result{
		Code: resp.Code,
		Msg:  resp.Msg,
	}
	if resp.Data != nil {
		result.Data = resp.Data.AsMap()
	}
	return result
}
