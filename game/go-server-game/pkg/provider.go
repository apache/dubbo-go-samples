package pkg

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/logger"

	"github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

type BasketballService struct{}

func (p *BasketballService) Online(ctx context.Context, uid string) (*pojo.Result, error) {
	logger.Infof("online: %s", uid)
	return &pojo.Result{Code: 0}, nil
}

func (p *BasketballService) Offline(ctx context.Context, uid string) (*pojo.Result, error) {
	logger.Infof("offline: %#s", uid)
	return &pojo.Result{Code: 0}, nil
}

func (p *BasketballService) Message(ctx context.Context, uid, data string) (*pojo.Result, error) {
	logger.Infof("message: %#v, %#v", uid, data)

	// // auto reply the same message
	// _, err := gateBasketball.Send(context.TODO(), uid, data)
	// if err != nil {
	//     logger.Errorf("send fail: %#s", err.Error())
	//     return &pojo.Result{Code: 1, Msg: err.Error()}, err
	// }

	return &pojo.Result{Code: 0, Data: map[string]interface{}{"to": uid, "message": data}}, nil
}

func (p *BasketballService) Reference() string {
	return "gameProvider.basketballService"
}

type JumpService struct{}

func (p *JumpService) Online(ctx context.Context, uid string) (*pojo.Result, error) {
	logger.Infof("online: %#s", uid)
	return &pojo.Result{Code: 0}, nil
}

func (p *JumpService) Offline(ctx context.Context, uid string) (*pojo.Result, error) {
	logger.Infof("offline: %#s", uid)
	return &pojo.Result{Code: 0}, nil
}

func (p *JumpService) Message(ctx context.Context, uid, data string) (*pojo.Result, error) {
	logger.Infof("message: %#s, %#s", uid, data)

	// reply the same message
	_, err := gateBasketball.Send(context.TODO(), uid, data)
	if err != nil {
		logger.Errorf("send fail: %#s", err.Error())
		return &pojo.Result{Code: 1, Msg: err.Error()}, err
	}

	return &pojo.Result{Code: 0, Data: map[string]interface{}{"to": uid, "message": data}}, nil
}

func (p *JumpService) Reference() string {
	return "gameProvider.jumpService"
}
