package pkg

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	hessian "github.com/apache/dubbo-go-hessian2"
)

func init() {
	config.SetProviderService(new(UserProvider))
	config.SetProviderService(new(UserProviderTriple))
	// ------for hessian2------
	hessian.RegisterPOJO(&User{})
	hessian.RegisterPOJO(&UserResponse{})
}
