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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
	"github.com/apache/dubbo-go-samples/game/go-server-gate/pkg"
	"github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

func init() {
	config.SetProviderService(new(pkg.BasketballService))

	config.SetConsumerService(pkg.GameBasketball)

	hessian.RegisterPOJO(&pojo.Result{})
}

func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	go startHttp()

	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			logger.Infof("app need reload")
		default:
			time.AfterFunc(time.Duration(time.Second*5), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			// The program exits normally or timeout forcibly exits.
			logger.Warnf("app exit now...")
			return
		}
	}
}

func startHttp() {

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		res, err := pkg.Login(context.TODO(), r.URL.Query().Get("name"))
		if err != nil {
			responseWithOrigin(w, r, 200, []byte(err.Error()))
			return
		}

		b, err := json.Marshal(res)
		if err != nil {
			responseWithOrigin(w, r, 200, []byte(err.Error()))
			return
		}
		responseWithOrigin(w, r, 200, b)
	})

	http.HandleFunc("/score", func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error(err.Error())
		}
		var info pojo.Info
		err = json.Unmarshal(reqBody, &info)
		if err != nil {
			responseWithOrigin(w, r, 500, []byte(err.Error()))
			return
		}
		res, err := pkg.Score(context.TODO(), info.Name, strconv.Itoa(info.Score))
		if err != nil {
			responseWithOrigin(w, r, 200, []byte(err.Error()))
			return
		}

		b, err := json.Marshal(res)
		if err != nil {
			responseWithOrigin(w, r, 200, []byte(err.Error()))
			return
		}
		responseWithOrigin(w, r, 200, b)
	})

	http.HandleFunc("/rank", func(w http.ResponseWriter, r *http.Request) {
		res, err := pkg.Rank(context.TODO(), r.URL.Query().Get("name"))
		if err != nil {
			responseWithOrigin(w, r, 200, []byte(err.Error()))
			return
		}
		b, err := json.Marshal(res)
		if err != nil {
			responseWithOrigin(w, r, 200, []byte(err.Error()))
			return
		}
		responseWithOrigin(w, r, 200, b)
	})

	_ = http.ListenAndServe("127.0.0.1:8089", nil)
}

// avoid cors
func responseWithOrigin(w http.ResponseWriter, r *http.Request, code int, json []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, err := w.Write(json)
	if err != nil {
		panic(err)
	}
}
