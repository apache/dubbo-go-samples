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
	"fmt"
	"strconv"
)

import (
	"github.com/dubbogo/gost/log/logger"

	"google.golang.org/protobuf/types/known/structpb"
)

import (
	game "github.com/apache/dubbo-go-samples/game/proto/game"
	gateProto "github.com/apache/dubbo-go-samples/game/proto/gate"
)

type Info struct {
	Name  string
	Score int
}

type GameServiceHandler struct{}

var userMap = make(map[string]*Info)

func (p *GameServiceHandler) Login(ctx context.Context, req *game.LoginRequest) (*game.Result, error) {
	logger.Infof("message: %#v", req.Uid)
	var (
		info *Info
		ok   bool
	)

	// call gate service
	rsp, err := GateFootball.Send(context.TODO(), &gateProto.SendRequest{
		Uid:  req.Uid,
		Data: "",
	})
	if err != nil {
		logger.Errorf("send fail: %#s", err.Error())
		return &game.Result{Code: 1, Msg: err.Error()}, err
	}

	fmt.Println("receive data from gate:", rsp)

	if info, ok = userMap[req.Uid]; !ok {
		info = &Info{}
		info.Name = req.Uid
		userMap[req.Uid] = info
	}

	data, _ := structpb.NewStruct(map[string]interface{}{
		"to":    req.Uid,
		"score": info.Score,
	})

	return &game.Result{
		Code: 0,
		Msg:  info.Name + ", your score is " + strconv.Itoa(info.Score),
		Data: data,
	}, nil
}

func (p *GameServiceHandler) Score(ctx context.Context, req *game.ScoreRequest) (*game.Result, error) {
	logger.Infof("message: %#v, %#v", req.Uid, req.Score)
	var (
		info *Info
		ok   bool
	)

	// call gate service
	rsp, err := GateFootball.Send(context.TODO(), &gateProto.SendRequest{
		Uid:  req.Uid,
		Data: req.Score,
	})
	if err != nil {
		logger.Errorf("send fail: %#s", err.Error())
		return &game.Result{Code: 1, Msg: err.Error()}, err
	}

	fmt.Println("receive data from gate:", rsp)

	if info, ok = userMap[req.Uid]; !ok {
		info = &Info{
			Name: req.Uid,
		}
		userMap[req.Uid] = info
		logger.Error("user data not found")
		data, _ := structpb.NewStruct(map[string]interface{}{})
		return &game.Result{Code: 1, Msg: "user data not found", Data: data}, nil
	}
	intSource, err := strconv.Atoi(req.Score)
	if err != nil {
		logger.Error(err.Error())
	}
	info.Score += intSource

	data, _ := structpb.NewStruct(map[string]interface{}{
		"to":    req.Uid,
		"score": info.Score,
	})

	return &game.Result{
		Code: 0,
		Msg:  "进球成功",
		Data: data,
	}, nil
}

func (p *GameServiceHandler) Rank(ctx context.Context, req *game.RankRequest) (*game.Result, error) {
	var (
		rank = 1
		info *Info
		ok   bool
	)

	// call gate service
	rsp, err := GateFootball.Send(context.TODO(), &gateProto.SendRequest{
		Uid:  req.Uid,
		Data: "",
	})
	if err != nil {
		logger.Errorf("send fail: %#s", err.Error())
		return &game.Result{Code: 1, Msg: err.Error()}, err
	}

	fmt.Println("receive data from gate:", rsp)

	if info, ok = userMap[req.Uid]; !ok {
		logger.Error("no user found")
		data, _ := structpb.NewStruct(map[string]interface{}{
			"to":   req.Uid,
			"rank": rank,
		})
		return &game.Result{Code: 1, Msg: "no user found", Data: data}, nil
	}

	for _, v := range userMap {
		if v.Score > info.Score {
			rank++
		}
	}

	data, _ := structpb.NewStruct(map[string]interface{}{
		"to":   req.Uid,
		"rank": rank,
	})

	return &game.Result{
		Code: 0,
		Msg:  "success",
		Data: data,
	}, nil
}
