package integration

import (
	"context"
	"testing"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"

	hessian "github.com/apache/dubbo-go-hessian2"
)

type UserProviderWithCustomGroupAndVersion struct {
	GetUser func(ctx context.Context, req *User) (rsp *User, err error)
}

type UserProvider struct {
	GetUser func(ctx context.Context, req *User) (rsp *User, err error)
}

type User struct {
	ID   string
	Name string
	Age  int32
	Time time.Time
}

func (u *User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

var userProvider = &UserProvider{}
var userProviderWithCustomRegistryGroupAndVersion = &UserProviderWithCustomGroupAndVersion{}

func TestMain(m *testing.M) {

	config.SetConsumerService(userProvider)
	config.SetConsumerService(userProviderWithCustomRegistryGroupAndVersion)
	hessian.RegisterPOJO(&User{})
	err := config.Load()
	if err != nil {
		panic(err)
	}

}
