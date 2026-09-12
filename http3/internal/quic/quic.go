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

// Package quic holds the QUIC transport options shared by the HTTP/3 sample server and client.
package quic

import (
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/protocol/triple"
)

// QUIC transport tuning values; unset options fall back to quic-go defaults.
const (
	keepAlivePeriod                = 30 * time.Second
	maxIdleTimeout                 = 90 * time.Second
	maxIncomingStreams             = 1024
	maxIncomingUniStreams          = 1024
	initialStreamReceiveWindow     = 512 * 1024
	maxStreamReceiveWindow         = 2 * 1024 * 1024
	initialConnectionReceiveWindow = 2 * 1024 * 1024
	maxConnectionReceiveWindow     = 8 * 1024 * 1024
)

// Options returns the HTTP/3 options applied by the sample server and client.
func Options() []triple.Option {
	return []triple.Option{
		triple.WithHttp3Enable(),
		triple.WithHttp3KeepAlivePeriod(keepAlivePeriod),
		triple.WithHttp3MaxIdleTimeout(maxIdleTimeout),
		triple.WithHttp3MaxIncomingStreams(maxIncomingStreams),
		triple.WithHttp3MaxIncomingUniStreams(maxIncomingUniStreams),
		triple.WithHttp3InitialStreamReceiveWindow(initialStreamReceiveWindow),
		triple.WithHttp3MaxStreamReceiveWindow(maxStreamReceiveWindow),
		triple.WithHttp3InitialConnectionReceiveWindow(initialConnectionReceiveWindow),
		triple.WithHttp3MaxConnectionReceiveWindow(maxConnectionReceiveWindow),
	}
}
