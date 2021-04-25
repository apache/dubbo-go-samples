package pkg

import (
    "context"
    "github.com/apache/dubbo-go/common/logger"

    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

type BasketballService struct{}

func (p *BasketballService) Send(ctx context.Context, uid, data string) (*pojo.Result, error) {
    logger.Infof("basketball: to=%s, message=%s", uid, data)
    return &pojo.Result{Code: 0, Data: map[string]interface{}{"to": uid, "message": data}}, nil
}

func (p *BasketballService) Reference() string {
    return "gateProvider.basketballService"
}
