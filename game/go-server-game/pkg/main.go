package pkg

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

func init() {
	config.SetProviderService(new(BasketballService))
	// config.SetProviderService(new(JumpService))

	config.SetConsumerService(gateBasketball)

	hessian.RegisterPOJO(&pojo.Result{})
}
