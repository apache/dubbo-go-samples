package pkg

import (
    "context"
    "github.com/apache/dubbo-go-samples/game/pkg/consumer/game"
    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

var GameBasketball = new(game.BasketballService)
// var gameJump = new(game.JumpService)


// just easy for demo test

func Message(ctx context.Context, uid string, data string) (*pojo.Result, error) {
    return GameBasketball.Message(ctx, uid, data)
    // return gameJump.Message(ctx, uid, data)
}

func Online(ctx context.Context, uid string) (*pojo.Result, error) {
    return GameBasketball.Online(ctx, uid)
    // return gameJump.Online(ctx, uid)
}

func Offline(ctx context.Context, uid string) (*pojo.Result, error) {
    return GameBasketball.Offline(ctx, uid)
    // return gameJump.Offline(ctx, uid)
}

func Login(ctx context.Context, data string) (*pojo.Result, error) {
    return GameBasketball.Login(ctx, data)
}

func Score(ctx context.Context, uid string, score int) (*pojo.Result, error) {
    return GameBasketball.Score(ctx, uid, score)
}
