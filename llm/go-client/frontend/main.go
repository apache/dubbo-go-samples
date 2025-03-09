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
)

import (
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

import (
	"github.com/apache/dubbo-go-samples/llm/go-client/frontend/handlers"
	"github.com/apache/dubbo-go-samples/llm/go-client/frontend/service"
	chat "github.com/apache/dubbo-go-samples/llm/proto"
)

func main() {
	// init Dubbo
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(fmt.Sprintf("Error creating Dubbo client: %v", err))
	}

	svc, err := chat.NewChatService(cli)
	if err != nil {
		panic(fmt.Sprintf("Error creating chat service: %v", err))
	}

	// init Gin
	r := gin.Default()

	// config session
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("llm_session", store))

	// register tmpl
	r.LoadHTMLGlob("go-client/frontend/templates/*")
	r.Static("../static", "go-client/frontend/static/")

	// init service
	ctxManager := service.NewContextManager()

	// register route
	h := handlers.NewChatHandler(svc, ctxManager)
	r.GET("/", h.Index)
	r.POST("/api/chat", h.Chat)
	r.POST("/api/context/new", h.NewContext)
	r.GET("/api/context/list", h.ListContexts)
	r.POST("/api/context/switch", h.SwitchContext)

	if err := r.Run(":8080"); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
