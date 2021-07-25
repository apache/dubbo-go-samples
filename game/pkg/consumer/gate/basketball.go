package gate

import (
	"context"
)

import (
	"github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

type BasketballService struct {
	Send func(ctx context.Context, uid string, data string) (*pojo.Result, error)
}

func (p *BasketballService) Reference() string {
	return "gateConsumer.basketballService"
}
