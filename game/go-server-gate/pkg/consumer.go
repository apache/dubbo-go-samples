package pkg

import (
    "context"
    "github.com/apache/dubbo-go-samples/game/pkg/consumer/game"
    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

var gameBasketball = new(game.BasketballService)
// var gameJump = new(game.JumpService)


// just easy for demo test

func Message(ctx context.Context, uid string, data string) (*pojo.Result, error) {
    return gameBasketball.Message(ctx, uid, data)
    // return gameJump.Message(ctx, uid, data)
}

func Online(ctx context.Context, uid string) (*pojo.Result, error) {
    return gameBasketball.Online(ctx, uid)
    // return gameJump.Online(ctx, uid)
}

func Offline(ctx context.Context, uid string) (*pojo.Result, error) {
    return gameBasketball.Offline(ctx, uid)
    // return gameJump.Offline(ctx, uid)
}
