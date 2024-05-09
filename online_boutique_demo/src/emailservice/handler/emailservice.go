package handler

import (
	"context"

	"go-micro.dev/v4/logger"

	pb "github.com/go-micro/demo/emailservice/proto"
)

type DummyEmailService struct{}

func (s *DummyEmailService) SendOrderConfirmation(ctx context.Context, in *pb.SendOrderConfirmationRequest, out *pb.Empty) error {
	logger.Infof("A request to send order confirmation email to %s has been received.", in.Email)
	return nil
}
