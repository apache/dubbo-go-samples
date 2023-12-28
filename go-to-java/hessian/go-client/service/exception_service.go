package service

import (
	"context"

	"dubbo/client/model"
)

type ExceptionService struct {
	// GetException is the client for the GetException service.
	Throwable func(ctx context.Context, req string) (model.Exception, error) `dubbo:"getThrowable"`
}

func (s *ExceptionService) Reference() string {
	return "ExceptionService"
}
