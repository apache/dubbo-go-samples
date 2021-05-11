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
