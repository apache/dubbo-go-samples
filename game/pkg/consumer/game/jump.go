package game

import (
    "context"
    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

type JumpService struct {
    Online  func(ctx context.Context, uid string) (*pojo.Result, error)
    Offline func(ctx context.Context, uid string) (*pojo.Result, error)
    Message func(ctx context.Context, uid string, data string) (*pojo.Result, error)
}

func (p *JumpService) Reference() string {
    return "gameConsumer.jumpService"
}