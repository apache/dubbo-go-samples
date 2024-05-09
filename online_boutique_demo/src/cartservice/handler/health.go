package handler

import (
	"context"

	pb "github.com/go-micro/demo/cartservice/proto"
)

type Health struct{}

func (h *Health) Check(ctx context.Context, req *pb.HealthCheckRequest, rsp *pb.HealthCheckResponse) error {
	rsp.Status = pb.HealthCheckResponse_SERVING
	return nil
}

func (h *Health) Watch(ctx context.Context, req *pb.HealthCheckRequest, stream pb.Health_WatchStream) error {
	stream.Send(&pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	})
	return nil
}
