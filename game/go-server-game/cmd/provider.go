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
	"strconv"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
)

import (
	"github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

type BasketballService struct{}

var userMap = make(map[string]*pojo.Info)

func (p *BasketballService) Login(ctx context.Context, uid string) (*pojo.Result, error) {
	logger.Infof("message: %#v", uid)
	var (
		info *pojo.Info
		ok   bool
	)

	// auto reply the same message
	rsp, err := GateBasketball.Send(context.TODO(), uid, "")
	if err != nil {
		logger.Errorf("send fail: %#s", err.Error())
		return &pojo.Result{Code: 1, Msg: err.Error()}, err
	}

	fmt.Println("receive data from gate:", rsp)

	if info, ok = userMap[uid]; !ok {
		info = &pojo.Info{}
		info.Name = uid
		userMap[uid] = info
	}
	return &pojo.Result{Code: 0, Msg: info.Name + ", your score is " + strconv.Itoa(info.Score), Data: map[string]interface{}{"to": uid, "score": info.Score}}, nil
}

func (p *BasketballService) Score(ctx context.Context, uid, score string) (*pojo.Result, error) {
	logger.Infof("message: %#v, %#v", uid, score)
	var (
		info = &pojo.Info{}
		ok   bool
	)

	// auto reply the same message
	rsp, err := GateBasketball.Send(context.TODO(), uid, score)
	if err != nil {
		logger.Errorf("send fail: %#s", err.Error())
		return &pojo.Result{Code: 1, Msg: err.Error()}, err
	}

	fmt.Println("receive data from gate:", rsp)

	if info, ok = userMap[uid]; !ok {
		info = &pojo.Info{
			Name: uid,
		}
		userMap[uid] = info
		logger.Error("user data not found")
		return &pojo.Result{Code: 1, Msg: "user data not found", Data: map[string]interface{}{}}, nil
	}
	intSource, err := strconv.Atoi(score)
	if err != nil {
		logger.Error(err.Error())
	}
	info.Score += intSource

	return &pojo.Result{Code: 0, Msg: "进球成功", Data: map[string]interface{}{"to": uid, "score": info.Score}}, nil
}

func (p *BasketballService) Rank(ctx context.Context, uid string) (*pojo.Result, error) {
	var (
		rank = 1
		info *pojo.Info
		ok   bool
	)

	// auto reply the same message
	rsp, err := GateBasketball.Send(context.TODO(), uid, "")
	if err != nil {
		logger.Errorf("send fail: %#s", err.Error())
		return &pojo.Result{Code: 1, Msg: err.Error()}, err
	}

	fmt.Println("receive data from gate:", rsp)

	if info, ok = userMap[uid]; !ok {
		logger.Error("no user found")
		return &pojo.Result{Code: 1, Msg: "no user found", Data: map[string]interface{}{"to": uid, "rank": rank}}, nil
	}

	for _, v := range userMap {
		if v.Score > info.Score {
			rank++
		}
	}

	return &pojo.Result{Code: 0, Msg: "success", Data: map[string]interface{}{"to": uid, "rank": rank}}, nil
}

func (p *BasketballService) Reference() string {
	return "gameProvider.basketballService"
}
