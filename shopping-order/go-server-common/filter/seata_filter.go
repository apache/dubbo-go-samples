package filter

import (
	"context"
)

import (
	"github.com/apache/dubbo-go/common/extension"
	"github.com/apache/dubbo-go/filter"
	"github.com/apache/dubbo-go/protocol"
)

import (
// inner
)

const (
	SEATA     = "seata"
	SEATA_XID = "seata_xid"
)

func init() {
	/**
	 * `seata` would be the name that used in your configuration file.
	 * it can be used as reference filter and provider filter.
	 *
	 * filter: "seata"
	 *
	 * registries:
	 *  "demoZk":
	 *    protocol: "zookeeper"
	 *    timeout	: "3s"
	 *    address: "127.0.0.1:2181"
	 *
	 * Another important things is that you should make sure this statement executed. It usually means that
	 * this file should be imported.
	 */
	extension.SetFilter(SEATA, getSeataFilter)
}

type SeataFilter struct{}

func (sf *SeataFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {

	// get transaction XID from invocation provides, or more information about this
	// maybe you can get many params in url : `url := invoker.GetUrl()`, `url.ServiceKey()`
	xid := invocation.AttachmentsByKey(SEATA_XID, "")
	if xid != "" {
		return invoker.Invoke(context.WithValue(ctx, SEATA_XID, xid), invocation)
	}
	return invoker.Invoke(ctx, invocation)
}

func (sf *SeataFilter) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	return result
}

func getSeataFilter() filter.Filter {
	return &SeataFilter{}
}
