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

package filter

import (
	"context"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
)

const (
	SEATA     = "seata"
	SEATA_XID = "seata_xid"
)

func init() {
	extension.SetFilter(SEATA, getSeataFilter)
}

// SeataFilter ...
type SeataFilter struct {
}

// Invoke ...
func (sf *SeataFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	xid := invocation.AttachmentsByKey(SEATA_XID, "")
	if xid != "" {
		return invoker.Invoke(context.WithValue(ctx, SEATA_XID, xid), invocation)
	}
	return invoker.Invoke(ctx, invocation)
}

// OnResponse ...
func (sf *SeataFilter) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	return result
}

// GetSeataFilter ...
func getSeataFilter() filter.Filter {
	return &SeataFilter{}
}
