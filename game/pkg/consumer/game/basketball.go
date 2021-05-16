package game

import (
    "context"
    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

type BasketballService struct {
    Login   func(ctx context.Context, uid string) (*pojo.Result, error)
    Score   func(ctx context.Context, uid, score string) (*pojo.Result, error)
    Rank  func(ctx context.Context, uid string) (*pojo.Result, error)
}

func (p *BasketballService) Reference() string {
    return "gameConsumer.basketballService"
}
