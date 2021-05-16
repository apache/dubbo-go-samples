package game

import (
    "context"
    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

type BasketballService struct {
    Online  func(ctx context.Context, uid string) (*pojo.Result, error)
    Offline func(ctx context.Context, uid string) (*pojo.Result, error)
    Message func(ctx context.Context, uid string, data string) (*pojo.Result, error)
    Login   func(ctx context.Context, data string) (*pojo.Result, error)
    Score   func(ctx context.Context, uid string, score int) (*pojo.Result, error)
}

func (p *BasketballService) Reference() string {
    return "gameConsumer.basketballService"
}
