package gate

import (
	"context"
	"github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

type JumpService struct {
	Send func(ctx context.Context, uid string, data string) (*pojo.Result, error)
}

func (p *JumpService) Reference() string {
	return "gateConsumer.jumpService"
}
