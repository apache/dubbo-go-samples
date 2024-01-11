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

package pages

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/dubbo-go-samples/task/shop/frontend/api"
	"github.com/apache/dubbo-go-samples/task/shop/frontend/server_v1"
	"github.com/gin-gonic/gin"
)

var (
	shopServer api.ShopService
)

func init() {
	provider, err := server_v1.NewShopServiceProvider()
	if err != nil {
		panic(err)
	}
	shopServer = provider
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func Login(c *gin.Context) {
	// get the query parameters
	username := c.Query("username")
	password := c.Query("password")
	// login
	if ok := shopServer.Login(username, password); !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "login failed",
		})
		return
	}
	//get item detail
	item, err := shopServer.CheckItem(1, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("get item failed error: %s", err.Error()),
		})
		return
	}
	// return html
	c.HTML(http.StatusOK, "detail.html", gin.H{"username": username, "item": item})
}

func TimeoutLogin(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{"result": "Failed to login, request timeout, please add timeout policy and retry!"})
		}
	}()
	// get the query parameters
	username := c.Query("username")
	password := c.Query("password")
	if ok := shopServer.TimeoutLogin(username, password); !ok {
		c.HTML(http.StatusOK, "index.html", gin.H{"result": "Failed to login, request timeout, please add timeout policy and retry!"})
		return
	}
	c.HTML(http.StatusOK, "detail.html", gin.H{"username": username})
}

func GrayLogin(c *gin.Context) {
	// get the query parameters
	username := c.Query("username")
	password := c.Query("password")
	// login
	if ok := shopServer.Login(username, password); !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "login failed",
		})
		return
	}
	//get item detail
	item, err := shopServer.CheckItemGray(1, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("get item failed error: %s", err.Error()),
		})
		return
	}
	// return html
	c.HTML(http.StatusOK, "detail.html", gin.H{"username": username, "item": item})
}

func UserInfo(c *gin.Context) {
	username := c.Query("username")
	user, err := shopServer.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("get item failed error: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateOrder(c *gin.Context) {
	// get the query parameters
	username := c.PostForm("username")
	sku := c.PostForm("sku")
	skuInt, _ := strconv.Atoi(sku)
	// request
	order, err := shopServer.SubmitOrder(int64(skuInt), 1, "beijing", "11111111", username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("get item failed error: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, order)
}
