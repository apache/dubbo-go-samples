package pkg

import (
    "context"
    "fmt"

    "github.com/apache/dubbo-go/common/logger"

    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

type BasketballService struct{}

func (p *BasketballService) Online(ctx context.Context, uid string) (*pojo.Result, error) {
    logger.Infof("online: %s", uid)
    return &pojo.Result{Code: 0, Msg: "hello this is game provider"}, nil
}

func (p *BasketballService) Offline(ctx context.Context, uid string) (*pojo.Result, error) {
    logger.Infof("offline: %#s", uid)
    return &pojo.Result{Code: 0, Msg: "hello this is game provider"}, nil
}

func (p *BasketballService) Message(ctx context.Context, uid, data string) (*pojo.Result, error) {
    logger.Infof("message: %#v, %#v", uid, data)

    // auto reply the same message
    rsp, err := GateBasketball.Send(context.TODO(), uid, data)
    if err != nil {
        logger.Errorf("send fail: %#s", err.Error())
        return &pojo.Result{Code: 1, Msg: err.Error()}, err
    }
    fmt.Println("receive data from gate:", rsp)
    return &pojo.Result{Code: 0, Data: map[string]interface{}{"to": uid, "message": data + " from game provider"}}, nil
}

func (p *BasketballService) Reference() string {
    return "gameProvider.basketballService"
}