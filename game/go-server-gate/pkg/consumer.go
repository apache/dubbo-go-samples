package pkg

import (
    "context"
    "github.com/apache/dubbo-go-samples/game/pkg/consumer/game"
    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

var GameBasketball = new(game.BasketballService)
// var gameJump = new(game.JumpService)


// just easy for demo test

func Login(ctx context.Context, data string) (*pojo.Result, error) {
    return GameBasketball.Login(ctx, data)
}

func Score(ctx context.Context, uid, score string) (*pojo.Result, error) {
    return GameBasketball.Score(ctx, uid, score)
}

func Rank (ctx context.Context, uid string) (*pojo.Result, error) {
    return GameBasketball.Rank(ctx, uid)
}
