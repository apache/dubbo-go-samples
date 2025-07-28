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
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	commentAPI "github.com/apache/dubbo-go-samples/task/shop/comment/api"
	"github.com/apache/dubbo-go-samples/task/shop/detail/api"
)

// DetailProvider is the provider of detail service
type DetailProvider struct {
	commentService commentAPI.Comment
}

func (d *DetailProvider) GetItem(ctx context.Context, req *api.GetItemReq) (*api.Item, error) {
	//get comment from comment server
	comment, err := d.commentService.GetComment(context.Background(), &commentAPI.CommentReq{
		ItemName: "wudong",
	})
	if err != nil {
		fmt.Printf("Detail provider get comment error: %v\n", err)
	}
	return &api.Item{
		Sku:         req.Sku,
		ItemName:    "shirt",
		Description: "item from detail v2",
		Stock:       100,
		Price:       100,
		Comment:     comment.Msg,
	}, nil
}

func (d *DetailProvider) DeductStock(ctx context.Context, req *api.DeductStockReq) (*api.DeductStockResp, error) {
	return &api.DeductStockResp{Success: true}, nil
}

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("shop-detail"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20021),
		),
	)
	if err != nil {
		panic(err)
	}

	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}
	commentService, err := commentAPI.NewComment(cli)
	if err != nil {
		panic(err)
	}
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	if err = api.RegisterDetailHandler(srv, &DetailProvider{commentService: commentService}); err != nil {
		panic(err)
	}
	if err = srv.Serve(); err != nil {
		logger.Error(err)
	}
}
