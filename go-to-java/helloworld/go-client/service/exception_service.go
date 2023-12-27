package service

import "context"

type ExceptionService struct {
	// GetException is the client for the GetException service.
	Throwable func(ctx context.Context, req string) error `dubbed:"getThrowable"`
}

func (s *ExceptionService) Reference() string {
	return "org.apache.dumbo.server.service.impl.ExceptionServiceImpl"
}
