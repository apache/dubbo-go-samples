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

package integration

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	game "github.com/apache/dubbo-go-samples/game/proto/game"
)

var gameService game.GameService

func TestMain(m *testing.M) {
	stopGate, err := ensureGateServer()
	if err != nil {
		panic(err)
	}

	cli, err := client.NewClient(
		client.WithClientURL(envOrDefault("GAME_SERVER_ADDR", "127.0.0.1:20000")),
	)
	if err != nil {
		panic(err)
	}

	gameService, err = game.NewGameService(cli)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	if stopGate != nil {
		stopGate()
	}

	os.Exit(code)
}

func envOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func ensureGateServer() (func(), error) {
	addr := envOrDefault("GATE_SERVER_ADDR", "127.0.0.1:20001")

	if portReachable(addr, 500*time.Millisecond) {
		return nil, nil
	}

	root, err := repoRoot()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("go", "run", "./game/gate/go-server/cmd")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start gate server: %w", err)
	}

	if err := waitForPort(addr, 15*time.Second); err != nil {
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("gate server failed to start: %w", err)
	}

	return func() {
		_ = cmd.Process.Signal(os.Interrupt)
		done := make(chan struct{})
		go func() {
			_ = cmd.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = cmd.Process.Kill()
		}
	}, nil
}

func waitForPort(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if portReachable(addr, 500*time.Millisecond) {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("timed out waiting for %s", addr)
}

func portReachable(addr string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", dir)
		}
		dir = parent
	}
}
