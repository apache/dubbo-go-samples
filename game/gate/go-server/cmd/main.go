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
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/game/gate/pkg"
	gateProto "github.com/apache/dubbo-go-samples/game/proto/gate"
)

type Info struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func main() {
	// Initialize game client
	cli, err := pkg.InitGameClient()
	if err != nil {
		logger.Fatalf("failed to create game client: %v", err)
	}
	pkg.SetGameClient(cli)

	// Start RPC server
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20001),
			protocol.WithTriple(),
		),
	)
	if err != nil {
		logger.Fatalf("failed to create server: %v", err)
	}

	if err := gateProto.RegisterGateServiceHandler(srv, &pkg.GateServiceHandler{}); err != nil {
		logger.Fatalf("failed to register gate service handler: %v", err)
	}

	// Start HTTP server in goroutine
	go startHttp()

	// Start signal handler in goroutine
	go initSignal()

	// Start RPC server
	if err := srv.Serve(); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}

func initSignal() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
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
	// Register API endpoints first (before static file server)
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
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err.Error())
		}
		var info Info
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

	// Serve static files from website directory (after API endpoints)
	websiteDir := resolveWebsiteDir()
	var fs http.Handler
	if websiteDir != "" {
		logger.Infof("serving static files from %s", websiteDir)
		fs = http.FileServer(http.Dir(websiteDir))
	} else {
		logger.Warnf("website directory not found, static files disabled")
		fs = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
	}
	// Only serve static files for non-API paths
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's an API endpoint
		if r.URL.Path == "/login" || r.URL.Path == "/score" || r.URL.Path == "/rank" {
			http.NotFound(w, r)
			return
		}
		fs.ServeHTTP(w, r)
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

func resolveWebsiteDir() string {
	if dir := os.Getenv("GAME_WEBSITE_DIR"); dir != "" {
		return dir
	}

	var roots []string
	if wd, err := os.Getwd(); err == nil {
		roots = append(roots, wd)
	}
	if exe, err := os.Executable(); err == nil {
		roots = append(roots, filepath.Dir(exe))
	}

	seen := make(map[string]struct{})
	for _, root := range roots {
		cur := root
		for i := 0; i < 6; i++ {
			if _, ok := seen[cur]; ok {
				break
			}
			seen[cur] = struct{}{}

			candidate := filepath.Join(cur, "website")
			if info, err := os.Stat(candidate); err == nil && info.IsDir() {
				return candidate
			}

			parent := filepath.Dir(cur)
			if parent == cur {
				break
			}
			cur = parent
		}
	}

	return ""
}
