package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/go-micro/demo/shippingservice/proto"
)

type Health struct{}

func (h *Health) Check(ctx context.Context, req *pb.HealthCheckRequest, rsp *pb.HealthCheckResponse) error {
	rsp.Status = pb.HealthCheckResponse_SERVING
	return nil
}

func (h *Health) Watch(ctx context.Context, req *pb.HealthCheckRequest, stream pb.Health_WatchStream) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}
