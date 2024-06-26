// Code generated by protoc-gen-triple. DO NOT EDIT.
//
// Source: payment.proto
package payment

import (
	"context"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
)

// This is a compile-time assertion to ensure that this generated file and the Triple package
// are compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of Triple newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of Triple or updating the Triple
// version compiled into your binary.
const _ = triple_protocol.IsAtLeastVersion0_1_0

const (
	// PaymentServiceName is the fully-qualified name of the PaymentService service.
	PaymentServiceName = "payment.PaymentService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// PaymentServiceChargeProcedure is the fully-qualified name of the PaymentService's Charge RPC.
	PaymentServiceChargeProcedure = "/payment.PaymentService/Charge"
)

var (
	_ PaymentService = (*PaymentServiceImpl)(nil)
)

// PaymentService is a client for the payment.PaymentService service.
type PaymentService interface {
	Charge(ctx context.Context, req *ChargeRequest, opts ...client.CallOption) (*ChargeResponse, error)
}

// NewPaymentService constructs a client for the payment.PaymentService service.
func NewPaymentService(cli *client.Client, opts ...client.ReferenceOption) (PaymentService, error) {
	conn, err := cli.DialWithInfo("payment.PaymentService", &PaymentService_ClientInfo, opts...)
	if err != nil {
		return nil, err
	}
	return &PaymentServiceImpl{
		conn: conn,
	}, nil
}

func SetConsumerService(srv common.RPCService) {
	dubbo.SetConsumerServiceWithInfo(srv, &PaymentService_ClientInfo)
}

// PaymentServiceImpl implements PaymentService.
type PaymentServiceImpl struct {
	conn *client.Connection
}

func (c *PaymentServiceImpl) Charge(ctx context.Context, req *ChargeRequest, opts ...client.CallOption) (*ChargeResponse, error) {
	resp := new(ChargeResponse)
	if err := c.conn.CallUnary(ctx, []interface{}{req}, resp, "Charge", opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

var PaymentService_ClientInfo = client.ClientInfo{
	InterfaceName: "payment.PaymentService",
	MethodNames:   []string{"Charge"},
	ConnectionInjectFunc: func(dubboCliRaw interface{}, conn *client.Connection) {
		dubboCli := dubboCliRaw.(*PaymentServiceImpl)
		dubboCli.conn = conn
	},
}

// PaymentServiceHandler is an implementation of the payment.PaymentService service.
type PaymentServiceHandler interface {
	Charge(context.Context, *ChargeRequest) (*ChargeResponse, error)
}

func RegisterPaymentServiceHandler(srv *server.Server, hdlr PaymentServiceHandler, opts ...server.ServiceOption) error {
	return srv.Register(hdlr, &PaymentService_ServiceInfo, opts...)
}

func SetProviderService(srv common.RPCService) {
	dubbo.SetProviderServiceWithInfo(srv, &PaymentService_ServiceInfo)
}

var PaymentService_ServiceInfo = server.ServiceInfo{
	InterfaceName: "payment.PaymentService",
	ServiceType:   (*PaymentServiceHandler)(nil),
	Methods: []server.MethodInfo{
		{
			Name: "Charge",
			Type: constant.CallUnary,
			ReqInitFunc: func() interface{} {
				return new(ChargeRequest)
			},
			MethodFunc: func(ctx context.Context, args []interface{}, handler interface{}) (interface{}, error) {
				req := args[0].(*ChargeRequest)
				res, err := handler.(PaymentServiceHandler).Charge(ctx, req)
				if err != nil {
					return nil, err
				}
				return triple_protocol.NewResponse(res), nil
			},
		},
	},
}
