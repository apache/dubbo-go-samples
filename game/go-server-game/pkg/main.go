package pkg

import (
    hessian "github.com/apache/dubbo-go-hessian2"
    "dubbo.apache.org/dubbo-go/v3/config"

    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

func init() {
    config.SetProviderService(new(BasketballService))
    // config.SetProviderService(new(JumpService))

    config.SetConsumerService(gateBasketball)

    hessian.RegisterPOJO(&pojo.Result{})
}
